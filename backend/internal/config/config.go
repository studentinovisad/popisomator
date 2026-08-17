package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	BackendPort uint16 `env:"BACKEND_PORT" envDefault:"8080"`
	JWTSecret   string `env:"POPISOMATOR_JWT_SECRET"`
	PostgresDSN string `env:"POPISOMATOR_POSTGRES_DSN,required"`
	DebugMode   bool   `env:"POPISOMATOR_DEBUG" envdefault:"false"`
}

var CurrentConfig Config
var InitDone = false

func (config Config) Address() string {
	return ":" + strconv.FormatUint(uint64(config.BackendPort), 10)
}

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
