package config

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host,omitempty"`
		Port     int    `yaml:"port,omitempty"`
		RootUser string `yaml:"root_user,omitempty"`
		RootPswd string `yaml:"root_pswd,omitempty"`
	} `yaml:"database,omitempty"`
}

var Cfg Config

func Read(slog *slog.Logger) error {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(configFile, Cfg); err != nil {
		return err
	}

	return nil
}
