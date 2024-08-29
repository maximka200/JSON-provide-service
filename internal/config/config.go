package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env     string        `mapstructure:"env"`
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
	Db      `mapstructure:"db"`
}

type Db struct {
	Username string `mapstructure:"username"`
	Password string
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DBName   string `mapstructure:"db_name"`
	Sslmode  bool   `mapstructure:"sslmode"`
}

func MustReadConfig() Config {
	var cfg Config

	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config/")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("unable to unmarshal into struct: %v", err))
	}

	return cfg
}
