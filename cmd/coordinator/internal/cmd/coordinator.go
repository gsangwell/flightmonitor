package cmd

import (
	"fmt"
	"os"
	"net/http"
	"log/slog"
	"flightmonitor/cmd/coordinator/internal/config"
	"flightmonitor/cmd/coordinator/internal/handlers"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "coordinator",
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
        rootCmd.PersistentFlags().StringVarP(&config.File, "config", "", "./coordinator.yaml", "Configuration file")

	// We do this here to force cobra to set flags
	rootCmd.Execute()

	// Add subcommands
	rootCmd.Short = "Start the Coordinator API."
        rootCmd.Long = `Start the Alces Monitoring Coordinator API`
	rootCmd.AddCommand(configureCmd)

	err := config.LoadConfig()

        if err != nil {
		// Ask for config
		rootCmd.Run = askForConfig
        } else {
		// Run the server
		rootCmd.Run = runServer
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}

func runServer(cmd *cobra.Command, args []string) {
	// API routes
        mux := http.NewServeMux()
        mux.HandleFunc("/status", handlers.Status)
        mux.HandleFunc("/register", handlers.Register)
        mux.HandleFunc("/registered", handlers.CheckRegistered)
        mux.HandleFunc("/managed/status", handlers.ManagedStatus)
        mux.HandleFunc("/managed/set", handlers.SetManaged)
        http.Handle("/", mux)

	slog.Info("Coordinator server starting on " + config.AppConfig.Coordinator.Listen + ":" + config.AppConfig.Coordinator.Port)
	http.ListenAndServe(config.AppConfig.Coordinator.Listen + ":" + config.AppConfig.Coordinator.Port, nil)
}

func askForConfig(cmd *cobra.Command, args []string) {
	// If we've ran the root command, and ended up here, we need to ask for config.
	fmt.Println("Please run 'coordinator configure' first.")
}
