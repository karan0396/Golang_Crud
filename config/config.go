package config

import (
	"api/pkg/logger"
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

// struct for config
type Config struct {
	Database struct {
		Driver string `json:"driver"`
		Dsn    string `json:"dsn"`
	} `json:"database"`

	Server struct {
		Port int `json:"port"`
	} `json:"server"`

	Pagination struct {
		Page  string `json:"page"`
		Limit string `json:"limit"`
	} `json:"pagination"`

	Delete struct {
		HardDelete int `json:"hardDelete"`
	} `json:"delete"`
}

func Load() (*Config, error) {
	logger.IntializeLogger()
	appConfig := &Config{}

	_, err := os.Stat("config.json")
	if err != nil {
		logger.Logger.DPanic("file not exit", zap.Error(err))
		return nil, err
	}

	file, err := os.Open("config.json")
	if err != nil {
		logger.Logger.DPanic("file not open", zap.Error(err))
		return nil, err
	}

	if err = json.NewDecoder(file).Decode(appConfig); err != nil {
		logger.Logger.DPanic("file not decode", zap.Error(err))
		return nil, err
	}
	return appConfig, nil
}
