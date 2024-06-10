package config

import (
	"main.go/internal/config/helper"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// HttpServer - contains ip and port for http server
type HttpServer struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type Config struct {
	HttpServer HttpServer `yaml:"http_server"`
}

func ReadConfigYaml(filePath string) (cfg Config, err error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return cfg, err
	}
	defer helper.Closer(file)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
