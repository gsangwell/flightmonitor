package config

import (
	"os"
	"fmt"
	"strings"
	"bufio"
        "regexp"
	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
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
	f, err := os.Open(File)

	if err != nil {
		return err
	}

	if f != nil {
		defer f.Close()
	}

	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&AppConfig)

	if err != nil {
		return err
	}

	return nil
}

func Configure(cmd *cobra.Command, args []string) {
	// coordinator:
	//  listen: 0.0.0.0
	//  port: 9001
	//server:
	//  protocol: http
	//  host: <roster-server>
	//  port: 8080
	//cluster:
	//  site: site
	//  name: cluster

	roster := readInputStrictMatches("Roster Server", `(?P<protocol>http|https)://(?P<host>.*):(?P<port>[0-9]*)`, 3)
	site := readInputStrict("Site Name", `(^[a-zA-Z0-9]*$)`)
	cluster := readInputStrict("Cluster Name", `(^[a-zA-Z0-9]*$)`)

	var new_config Config

	new_config.Coordinator.Listen = "0.0.0.0"
	new_config.Coordinator.Port = "9001"

        new_config.Server.Protocol = roster[1]
        new_config.Server.Host = roster[2]
        new_config.Server.Port = roster[3]

	new_config.Cluster.Site = site
	new_config.Cluster.Name = cluster

        f, err := os.OpenFile(File, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

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

func readInputStrictMatches(prompt string, pattern string, expected int) []string {
        reader := bufio.NewReader(os.Stdin)

        var matches []string
        valid := false

        for !valid {
		fmt.Print(prompt + ": ")
		line, _ := reader.ReadString('\n')
		line, _ = strings.CutSuffix(line, "\n")
		line, _ = strings.CutSuffix(line, "\r")

                r := regexp.MustCompile(pattern)
                matches = r.FindStringSubmatch(line)

                if len(matches) == (expected  + 1){
                        valid = true
                }
        }

        return matches
}

func readInputStrict(prompt string, pattern string) string {
        reader := bufio.NewReader(os.Stdin)

        var matches []string
        valid := false

        for !valid {
                fmt.Print(prompt + ": ")
                line, _ := reader.ReadString('\n')
                line, _ = strings.CutSuffix(line, "\n")
                line, _ = strings.CutSuffix(line, "\r")

                r := regexp.MustCompile(pattern)
                matches = r.FindStringSubmatch(line)

                if len(matches) == 2 {
                        valid = true
                }
        }

        return matches[1]
}

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(prompt + ": ")
	line, _ := reader.ReadString('\n')
        line, _ = strings.CutSuffix(line, "\n")
        line, _ = strings.CutSuffix(line, "\r")

	return line
}

var AppConfig Config
var File string
