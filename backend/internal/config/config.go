package config

import (
	"errors"
	"os"
)

type Config struct {
	Address string
}

var CurrentConfig Config
var InitDone = false

func Init() error {
	if InitDone {
		return errors.New("Already initialised config")
	}

	if address_env, is_env := os.LookupEnv("POPISOMATOR_BACKEND_ADDR"); is_env {
		CurrentConfig.Address = address_env
	} else {
		CurrentConfig.Address = "localhost:8080"
	}

	InitDone = true
	return nil
}
