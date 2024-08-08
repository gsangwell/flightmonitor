package config

import (
        "os"
        "log/slog"
        "gopkg.in/yaml.v2"
)

type Config struct {
        Server struct {
                Protocol string `yaml:"protocol"`
                Port string `yaml:"port"`
                Host string `yaml:"host"`
        } `yaml:"server"`
}

func LoadConfig() error {
        f, err := os.Open("monitor.yaml")

        if err != nil {
                slog.Error("loadConfig", "err", err)
                return err
        }

        if f != nil {
                defer f.Close()
        }

        decoder := yaml.NewDecoder(f)

        err = decoder.Decode(&CliConfig)

        if err != nil {
                slog.Error("loadConfig", "f", f, "err", err)
                return err
        }

        return nil
}

var CliConfig Config
