package main

import (
	yaml "gopkg.in/yaml.v3"

	"github.com/spf13/viper"
)

const (
	DbTypeMemory     = "memory"
	DbTypePostgresql = "postgresql"
)

type Config struct {
	Env        string     `mapstructure:"env"`
	Logger     LoggerConf `mapstructure:"logger"`
	Db         DbConf     `mapstructure:"db"`
	HttpServer HttpServer `mapstructure:"http_server"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type DbConf struct {
	Type            string `mapstructure:"type"`
	Dsn             string `mapstructure:"dsn"`
	Migrations      string `mapstructure:"migrations"`
	MaxIdleConnects int    `mapstructure:"max_idle_connects"`
	MaxOpenConnects int    `mapstructure:"max_open_connects"`
}

type HttpServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewConfig() *Config {
	cfg := new(Config)
	return cfg
}

func (c *Config) initDefaults() error {
	viper.SetDefault("logger.level", "DEBUG")
	viper.SetDefault("logger.file", "")
	viper.SetDefault("db.type", DbTypeMemory)
	viper.SetDefault("db.dsn", "")
	viper.SetDefault("db.migrations", "")
	viper.SetDefault("db.max_idle_connects", "10")
	viper.SetDefault("db.max_open_connects", "10")
	viper.SetDefault("http_server.host", "0.0.0.0")
	viper.SetDefault("http_server.port", "8080")
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

func (c *Config) dump() (string, error) {
	s := viper.AllSettings()
	b, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
