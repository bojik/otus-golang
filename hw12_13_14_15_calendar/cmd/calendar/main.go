package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/api"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	internalapi "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/server/api"
	internalhttp "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	migrate "github.com/golang-migrate/migrate/v4"
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
	defer func() {
		_ = lf.Close()
	}()

	var (
		db       *sqlstorage.Storage
		calendar *app.App
	)
	if config.Db.Type == DbTypePostgresql {
		db = createDbStorage(config, migrateDown, fixAndForceMigration, logg)
		if db == nil {
			return
		}
		calendar = app.New(logg, db)
	} else {
		storage := memorystorage.New()
		calendar = app.New(logg, storage)
	}

	server := internalhttp.NewServer(logg, calendar, net.JoinHostPort(config.HttpServer.Host, config.HttpServer.Port))
	apiSever := internalapi.NewServer(
		net.JoinHostPort(config.ApiServer.Host, config.ApiServer.Port),
		api.NewCalendarApi(calendar, logg),
		logg,
	)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if config.Db.Type == DbTypePostgresql {
			if err := db.Close(ctx); err != nil {
				logg.Error("failed to close db connects: " + err.Error())
			}
		}

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		if err := apiSever.Stop(ctx); err != nil {
			logg.Error("failed to stop api server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := apiSever.Start(ctx); err != nil {
			logg.Error("failed to start api server: " + err.Error())
		}
	}()
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func createDbStorage(config *Config, migrateDown *int, fixAndForceMigration *int, logg logger.Logger) *sqlstorage.Storage {
	db := sqlstorage.New(config.Db.Dsn, config.Db.MaxIdleConnects, config.Db.MaxOpenConnects)

	if *migrateDown > 0 {
		logg.Info("migration down", logger.NewIntParam("steps", *migrateDown))
		if err := db.MigrateDown(config.Db.Migrations, *migrateDown); err != nil {
			logg.Error("migration down: " + err.Error())
			panic(err)
		}
		logg.Info("migration down: done")
		return nil
	}

	if *fixAndForceMigration > 0 {
		logg.Info("fix and force migration", logger.NewIntParam("version", *fixAndForceMigration))
		if err := db.FixAndForce(config.Db.Migrations, *fixAndForceMigration); err != nil {
			logg.Error("fix and force migration: " + err.Error())
			panic(err)
		}
		logg.Info("fix and force migration: done")
		return nil
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
	return db
}
