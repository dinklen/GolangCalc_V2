package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Microservice struct {
		Host           string        `mapstructure:"host"`
		Port           string        `mapstructure:"port"`
		ComputingPower int           `mapstructure:"computing_power"`
		WaitingTime    time.Duration `mapstructure:"waiting_time"`
		PlusTime       time.Duration `mapstructure:"plus_time"`
		MinusTime      time.Duration `mapstructure:"minus_time"`
		MultiplyTime   time.Duration `mapstructure:"multiply_time"`
		DivisionTime   time.Duration `mapstructure:"division_time"`
	} `mapstructure:"microservice"`

	JWT struct {
		AccessExpiry  time.Duration `mapstructure:"access_expiry"`
		RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
		Secret        string        `mapstructure:"secret"`
	} `mapstructure:"jwt"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"ssl_mode"`
	} `mapstructure:"database"`

	Redis struct {
		Host        string        `mapstructure:"host"`
		Port        string        `mapstructure:"port"`
		Password    string        `mapstructure:"password"`
		WaitingTime time.Duration `mapstructure:"waiting_time"`
	} `mapstructure:"redis"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../internal/config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
