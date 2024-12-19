package config

import (
	"github.com/joho/godotenv"
)

func Init() error {
	// read the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
