package config

import (
	"errors"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Address     string `env:"POPISOMATOR_BACKEND_ADDR" envDefault:"localhost:8080"`
	PostgresDSN string `env:"POPISOMATOR_POSTGRES_DSN,required"`
}

var CurrentConfig Config
var InitDone = false

func Init() error {
	if InitDone {
		return errors.New("Already initialised config")
	}

	for _, path := range []string{".env", "../.env"} {
		if err := godotenv.Load(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	if err := env.Parse(&CurrentConfig); err != nil {
		return err
	}

	InitDone = true
	return nil
}
