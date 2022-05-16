package main

import (
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"
)

const (
	DBTypePostgresql = config.DBTypePostgresql
)

type Config struct {
	Env        string            `mapstructure:"env"`
	Logger     config.LoggerConf `mapstructure:"logger"`
	DB         config.DBConf     `mapstructure:"db"`
	HTTPServer HTTPServer        `mapstructure:"http_server"`
	APIServer  APIServer         `mapstructure:"api_server"`
}

type HTTPServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type APIServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func NewConfig() *Config {
	cfg := new(Config)
	return cfg
}

func (c *Config) initDefaults() error {
	config.InitLoggerConfig()
	config.InitDBConfig()
	viper.SetDefault("http_server.host", "0.0.0.0")
	viper.SetDefault("http_server.port", "8080")
	viper.SetDefault("api_server.host", "0.0.0.0")
	viper.SetDefault("api_server.port", "8081")
	if err := viper.Unmarshal(&c); err != nil {
		return err
	}
	config.LoadFromEnv()
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
