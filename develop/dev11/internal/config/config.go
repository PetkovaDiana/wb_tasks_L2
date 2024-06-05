package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"main.go/server"
	"os"
)

const (
	envPath = ".env"
)

type App struct {
	ServerConfig *server.Config `yaml:"server"`
}

func NewAppConfig() (*App, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	cfgApp := new(App)

	if err := cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), cfgApp); err != nil {
		return nil, err
	}

	return cfgApp, nil
}
