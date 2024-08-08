package config

import (
        "fmt"
        "os"
        "bufio"
        "regexp"
        "strings"
        "github.com/spf13/cobra"
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
        f, err := os.Open(File)

        if err != nil {
                return err
        }

        if f != nil {
                defer f.Close()
        }

        decoder := yaml.NewDecoder(f)

        err = decoder.Decode(&CliConfig)

        if err != nil {
                return err
        }

        return nil
}

func Configure(cmd *cobra.Command, args []string) {
        reader := bufio.NewReader(os.Stdin)

        var matches []string
        valid := false

        for !valid {
                fmt.Print("Coordinator Server: ")
                server, _ := reader.ReadString('\n')
                server, _ = strings.CutSuffix(server, "\n")
                server, _ = strings.CutSuffix(server, "\r")

                r := regexp.MustCompile(`(?P<protocol>http|https)://(?P<host>.*):(?P<port>[0-9]*)`)
                matches = r.FindStringSubmatch(server)

                if len(matches) == 4 {
                        valid = true
                }
        }

        var new_config Config

        new_config.Server.Protocol = matches[1]
        new_config.Server.Host = matches[2]
        new_config.Server.Port = matches[3]

        f, err := os.OpenFile("./monitor.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

        if err != nil {
                fmt.Println("Error saving configuration.")
                return
        }

        if f != nil {
                defer f.Close()
        }

        enc := yaml.NewEncoder(f)

        err = enc.Encode(new_config)

        if err != nil {
                fmt.Println("Error saving configuration.")
                return
        }

        fmt.Println("Configuration saved!")
}

var CliConfig Config
var File string
