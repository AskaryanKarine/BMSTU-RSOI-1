package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppEnv      string `env:"APP_ENV" envDefault:"test"`
	PostgresDSN string `env:"POSTGRES_DSN"`
	Port        int    `env:"PORT" envDefault:"8000"`
}

func New() (Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil

}
