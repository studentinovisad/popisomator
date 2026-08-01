package config

import (
	"errors"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Address string `env:"POPISOMATOR_BACKEND_ADDR" envDefault:"localhost:8080"`
}

var CurrentConfig Config
var InitDone = false

func Init() error {
	if InitDone {
		return errors.New("Already initialised config")
	}

	if err := env.Parse(&CurrentConfig); err != nil {
		return err
	}

	InitDone = true
	return nil
}