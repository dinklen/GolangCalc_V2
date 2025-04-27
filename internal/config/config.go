package config

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	}

	Microservice struct {
		Port int `mapstructure:"port"`
	}

	Session struct {
		Id     int  `mapstructure:"id"`
		Status bool `mapstructure:"status"`
	}

	Database struct {
		Password string `mapstructure:"password"`
	}
}

func (cfg *Config) Load() *Config {
	return nil
}
