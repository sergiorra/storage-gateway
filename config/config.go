package config

import (
	"encoding/json"
	"os"

	"storage-gateway/internal/log"
)

type Config struct {
	App  App
	Api  Api
	Http Http
}

type App struct {
	LogLevel                 string
	ShutdownTimeoutInSeconds int
}

type Api struct {
	Host                       string
	Port                       int
	TimeoutInSeconds           int
	ReadHeaderTimeoutInSeconds int
}

type Http struct {
	MaxIdleConns        int
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	TimeoutInSeconds    int
}

func Read(filename string) (*Config, error) {
	var config Config

	log.Infof("Loading configuration from [%s]", filename)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
