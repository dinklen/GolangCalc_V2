package config

import (
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/service/errors/cfgerr"

	"github.com/spf13/viper"
	"go.uber.org/zap"
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
		PowerTime      time.Duration `mapstructure:"power_time"`
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

func Load(logger *zap.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../internal/config")
	logger.Info("config load success")

	if err := viper.ReadInConfig(); err != nil {
		return nil, cfgerr.ErrConfigReadingFailed
	}
	logger.Info("config read success")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, cfgerr.ErrConfigDataUnmarshalingFailed
	}
	logger.Info("config unmarshal success")

	return &cfg, nil
}

func Update(arg string, value any, logger *zap.Logger) error {
	viper.Set(arg, value)

	if err := viper.WriteConfig(); err != nil {
		return cfgerr.ErrConfigUpdattingFailed
	}

	logger.Info("config update success")
	return nil
}
