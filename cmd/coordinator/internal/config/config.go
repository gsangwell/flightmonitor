package config

import (
	"os"
	"log/slog"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Coordinator struct {
		Listen string `yaml:"listen"`
		Port string `yaml:"port"`
	} `yaml:"coordinator"`
	Server struct {
		Protocol string `yaml:"protocol"`
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Cluster struct  {
		Site string `yaml:"site"`
		Name string `yaml:"name"`
	} `yaml:"cluster"`
}

func LoadConfig() error {
	f, err := os.Open("coordinator.yaml")

	if err != nil {
		slog.Error("loadConfig", "err", err)
		return err
	}

	if f != nil {
		defer f.Close()
	}

	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&AppConfig)

	if err != nil {
		slog.Error("loadConfig", "f", f, "err", err)
		return err
	}

	return nil
}

var AppConfig Config
