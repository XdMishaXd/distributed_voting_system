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
	Clickhouse `yaml:"clickhouse"`
}

type Swagger struct {
	Username string `yaml:"username" env:"SWAGGER_USERNAME" env-default:"admin"`
	Password string `yaml:"password" env:"SWAGGER_PASSWORD" env-default:"admin"`
	Enabled  bool   `yaml:"enabled" env-default:"false"`
}

type Clickhouse struct {
	Addr     string `yaml:"addr" env-required:"true"`
	Database string `yaml:"database" env-required:"true"`
	Username string `yaml:"-" env:"CLICKHOUSE_USER" env-required:"true"`
	Password string `yaml:"-" env:"CLICKHOUSE_PASSWORD" env-required:"true"`

	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"1h"`
	DialTimeout     time.Duration `yaml:"dial_timeout" env-default:"5s"`

	BatchSize     int           `yaml:"batch_size" env-default:"2000"`
	FlushInterval time.Duration `yaml:"flush_interval" env-default:"1s"`
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
