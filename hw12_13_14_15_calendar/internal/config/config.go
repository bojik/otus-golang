package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	DBTypeMemory     = "memory"
	DBTypePostgresql = "postgresql"
)

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type DBConf struct {
	Type            string `mapstructure:"type"`
	Dsn             string `mapstructure:"dsn"`
	Migrations      string `mapstructure:"migrations"`
	MaxIdleConnects int    `mapstructure:"max_idle_connects"`
	MaxOpenConnects int    `mapstructure:"max_open_connects"`
}

type AMQPConf struct {
	URL          string `mapstructure:"url"`
	ExchangeName string `mapstructure:"exchange_name"`
	QueueName    string `mapstructure:"queue_name"`
}

func InitLoggerConfig() {
	viper.SetDefault("logger.level", "DEBUG")
	viper.SetDefault("logger.file", "")
}

func InitDBConfig() {
	viper.SetDefault("db.type", DBTypePostgresql)
	viper.SetDefault("db.dsn", "")
	viper.SetDefault("db.migrations", "")
	viper.SetDefault("db.max_idle_connects", "10")
	viper.SetDefault("db.max_open_connects", "10")
}

func InitAMQPConfig() {
	viper.SetDefault("amqp.url", "")
	viper.SetDefault("amqp.exchange_name", "")
	viper.SetDefault("amqp.queue_name", "")
}

func LoadFromEnv() {
	viper.SetEnvPrefix("calendar")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
