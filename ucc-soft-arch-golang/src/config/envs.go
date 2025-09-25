package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Envs interface {
	Get(key string) string
}

type envsImpl struct{}

func (envsImpl) Get(key string) string {
	return os.Getenv(key)
}

// LoadEnvs loads environment variables from an optional .env file.
// If the file is missing, it will NOT panic; environment variables
// can be provided via the environment (e.g. docker-compose). Only
// unexpected errors will cause a panic.
func LoadEnvs(filename ...string) Envs {
	err := godotenv.Load(filename...)
	if err != nil {
		// If the error is due to file not existing, ignore it.
		if errors.Is(err, os.ErrNotExist) {
			return &envsImpl{}
		}
		// godotenv may return a wrapped *os.PathError; check string fallback
		if pathErr, ok := err.(*os.PathError); ok {
			if errors.Is(pathErr.Err, os.ErrNotExist) {
				return &envsImpl{}
			}
		}
		// otherwise panic to surface unexpected errors
		panic("Error loading .env file: " + err.Error())
	}
	return &envsImpl{}
}
