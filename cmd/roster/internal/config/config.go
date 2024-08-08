package config

import (
	"os"
	"log/slog"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Roster struct {
		Listen string `yaml:"listen"`
		Port string `yaml:"port"`
	} `yaml:"roster"`
	Database struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	Slack struct {
		Token string `yaml:"token"`
		ChannelID string `yaml:"channel"`
	} `yaml:"slack"`
	Zabbix struct {
		API string `yaml:"api"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"zabbix"`
}

func LoadConfig() error {
	f, err := os.Open("roster.yaml")

	if err != nil {
		slog.Error("loadConfig", "err", err)
		return err
	}

	if f != nil {
		defer f.Close()
	}

	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&RosterConfig)

	if err != nil {
		slog.Error("loadConfig", "f", f, "err", err)
		return err
	}

	return nil
}

var RosterConfig Config
