package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const configPath = "config.yaml"

type conf struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host,omitempty"`
		Port     int    `yaml:"port,omitempty"`
		RootUser string `yaml:"root_user,omitempty"`
		RootPswd string `yaml:"root_pswd,omitempty"`
	} `yaml:"database,omitempty"`
}

func Get() (*conf, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg conf
	if err = yaml.Unmarshal(configFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
