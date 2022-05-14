package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/queue"
	sqlstorage "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	flag "github.com/spf13/pflag"
)

func main() {
	configFile := flag.StringP("config", "c", "configs/config_scheduler.yaml", "Path to configuration file")
	version := flag.BoolP("version", "", false, "Print version of application")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	if err := config.load(*configFile); err != nil {
		panic(err)
	}

	minLevelOpt, err := logger.NewOptionMinLevel(config.Logger.Level)
	if err != nil {
		panic(err)
	}
	logg := logger.New(minLevelOpt)
	lf, err := logg.AddLogFile(config.Logger.File)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = lf.Close()
	}()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	db := sqlstorage.New(config.DB.Dsn, config.DB.MaxIdleConnects, config.DB.MaxOpenConnects)
	if err := db.Connect(ctx); err != nil {
		panic(err)
	}
	defer db.Close(context.Background())

	calendar := app.New(logg, db)

	que := queue.New(config.AMQP.URL, config.AMQP.ExchangeName, "direct", config.AMQP.QueueName, logg)
	if err := que.Connect(); err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		logg.Info("calendar scheduler terminating...")
		que.Close()
		defer cancel()
	}()
	d, err := time.ParseDuration(config.SelectInterval)
	if err != nil {
		panic(err)
	}
	logg.Info("calendar scheduler starting...")
	if err := que.Publish(ctx, d, calendar); err != nil {
		cancel()
		panic(err)
	}
	<-ctx.Done()
}
