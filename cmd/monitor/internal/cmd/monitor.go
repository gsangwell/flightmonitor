package cmd

import (
        "fmt"
	"os"
        "flightmonitor/internal/common"
	"flightmonitor/cmd/monitor/internal/config"
	"flightmonitor/cmd/monitor/internal/coordinator"
        "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "monitor",
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Alces Monitoring CLI Version",
  Long:  `Show the version of Alces Monitoring CLI`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(common.VERSION)
  },
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&config.File, "config", "", "./monitor.yaml", "Configuration file")
        rootCmd.PersistentFlags().BoolVarP(&common.Debug, "debug", "d", false, "debug output")

	rootCmd.Execute()

	// Load config
	err := config.LoadConfig()

        if err != nil {
		rootCmd.Run = config.Configure
		return
        }

	// Check if node is registered
        registered, err := coordinator.CheckRegistered()

        if err != nil {
                fmt.Println(err.Error())
                return
        }

	if !registered {
                rootCmd.AddCommand(registerCmd)
		return
	}

	// Get status via coordinator
        services, err := coordinator.GetManagedStatus()

        if err != nil {
		fmt.Println("Error - unable to lookup current status")
		return
	}

	if (*services)["change-control"].Enabled == "0" {
		rootCmd.AddCommand(changeControlCmd)
        }

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(enableCmd)
	rootCmd.AddCommand(disableCmd)
}

func initConfig() {
        fmt.Println("Loading config!")
        // Load config
        err := config.LoadConfig()

        if err != nil {
                os.Exit(1)
        }
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
