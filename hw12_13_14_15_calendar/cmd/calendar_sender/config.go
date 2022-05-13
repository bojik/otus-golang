package main

import (
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/spf13/viper"
)

type Config struct {
	AMQP    config.AMQPConf   `mapstructure:"amqp"`
	Logger  config.LoggerConf `mapstructure:"logger"`
	Threads int               `mapstructure:"threads"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.initDefaults(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) initDefaults() error {
	config.InitLoggerConfig()
	config.InitAMQPConfig()
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	return nil
}

func (c *Config) load(file string) error {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	return nil
}
