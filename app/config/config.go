package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

type Server struct {
	AppPort         string        `env:"APP_PORT" envDefault:"7000"`
	AppHost         string        `env:"APP_HOST" envDefault:"localhost"`
	AppReadTimeout  time.Duration `env:"APP_READ_TIMEOUT" envDefault:"60s"`
	AppWriteTimeout time.Duration `env:"APP_WRITE_TIMEOUT" envDefault:"60s"`
	AppIdleTimeout  time.Duration `env:"APP_IDLE_TIMEOUT" envDefault:"60s"`
	AppAddress      string
}

type Postgres struct {
	Host     string `env:"DB_HOST,notEmpty"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	Database string `env:"DB_NAME,notEmpty"`
	User     string `env:"DB_USER,notEmpty"`
	Password string `env:"DB_PASSWORD,notEmpty"`
	SSLmode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

type Config struct {
	Server Server
	DB     Postgres
	Auth   AuthConfig
}

type AuthConfig struct {
	Salt           string        `env:"APP_SALT,notEmpty"`
	SigningKey     string        `env:"SIGNING_KEY,notEmpty"`
	AccessTokenTTL time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"12h"`
}

func InitConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while parsing .env: %w", err)
	}

	cfg.Server.AppAddress = ":" + cfg.Server.AppPort

	return cfg, nil
}
