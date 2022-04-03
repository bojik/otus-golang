package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/golang-migrate/migrate/v4"
	flag "github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

func main() {
	configFile := flag.StringP("config", "c", "configs/config.yaml", "Path to configuration file")
	dumpConfig := flag.BoolP("dump-config", "d", false, "Dump config with default values")
	migrateDown := flag.IntP("migrate-down", "", 0, "Step back of db migration")
	fixAndForceMigration := flag.IntP("fix-force", "", 0, "Force sets a migration version. It resets the dirty state to false.")
	version := flag.BoolP("version", "", false, "Print version of application")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	config := NewConfig()
	err := config.initDefaults()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%#v", config)
	if *dumpConfig {
		str, err := config.dump()
		if err != nil {
			panic(err)
		}
		fmt.Println(str)
		return
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
	defer lf.Close()

	db := sqlstorage.New(config.Db.Dsn, config.Db.MaxIdleConnects, config.Db.MaxOpenConnects)

	if *migrateDown > 0 {
		logg.Info("migration down", logger.NewIntParam("steps", *migrateDown))
		if err := db.MigrateDown(config.Db.Migrations, *migrateDown); err != nil {
			logg.Error("migration down: " + err.Error())
			panic(err)
		}
		logg.Info("migration down: done")
		return
	}

	if *fixAndForceMigration > 0 {
		logg.Info("fix and force migration", logger.NewIntParam("version", *fixAndForceMigration))
		if err := db.FixAndForce(config.Db.Migrations, *fixAndForceMigration); err != nil {
			logg.Error("fix and force migration: " + err.Error())
			panic(err)
		}
		logg.Info("fix and force migration: done")
		return
	}

	logg.Info("executing migrations")
	if err := db.Migrate(config.Db.Migrations); err != nil {
		if xerrors.Is(err, migrate.ErrNoChange) {
			logg.Info("migration: " + err.Error())
		} else {
			logg.Error("failed to execute db migrations: " + err.Error())
			panic(err)
		}
	}

	if err := db.Connect(context.Background()); err != nil {
		logg.Error("failed to connect to db: " + err.Error())
		panic(err)
	}

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := db.Close(ctx); err != nil {
			logg.Error("failed to close db connects: " + err.Error())
		}

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
