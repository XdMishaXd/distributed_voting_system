package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Swagger    `yaml:"swagger"`
}

type Swagger struct {
	Username string `yaml:"username" env:"SWAGGER_USERNAME" env-default:"admin"`
	Password string `yaml:"password" env:"SWAGGER_PASSWORD" env-default:"admin"`
	Enabled  bool   `yaml:"enabled" env-default:"false"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	HandlersTimeout time.Duration `yaml:"handlers_timeout" env-default:"5s"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	return &cfg
}
