package main

import (
	"fmt"
	"strings"
	"flightmonitor/internal/common"
	"flightmonitor/cmd/monitor/internal/config"
	"flightmonitor/cmd/monitor/internal/commands"
	"flightmonitor/cmd/monitor/internal/coordinator"
	"github.com/spf13/cobra"
)

var valid_services = []string{"alerting", "change-control"}

var rootCmd = &cobra.Command{
    Use:   "monitor",
    Short: "Alces Monitoring CLI",
    Long:  `Enable and disable components of the Alces Monitoring Service.`,
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Alces Monitoring CLI Version",
  Long:  `Show the version of Alces Monitoring CLI`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(common.VERSION)
  },
}

var statusCmd = &cobra.Command{
	Use:	"status <" + strings.Join(valid_services, "|") + ">",
	Short:	"Service Status",
	Long:	`Show whether a particular monitoring service is enabled or disabled on this host.`,
	Args:	cobra.ExactValidArgs(1),
	ValidArgs: valid_services,
	Run:	commands.StatusService,
}

var enableCmd = &cobra.Command{
        Use:    "enable <" + strings.Join(valid_services, "|") + ">",
        Short:  "Enable Service",
        Long:   `Enable a specific monitoring service on this host.`,
        Args:   cobra.ExactValidArgs(1),
        ValidArgs: valid_services,
        Run:    commands.EnableService,
}

var disableCmd = &cobra.Command{
        Use:    "disable <" + strings.Join(valid_services, "|") + ">",
        Short:  "Disable Service",
        Long:   `Disable a specific monitoring service on this host.`,
        Args:   cobra.ExactValidArgs(1),
        ValidArgs: valid_services,
        Run:    commands.DisableService,
}

var changeControlCmd = &cobra.Command{
	Use:    "change-control",
	Short:  "Change Control",
	Long:	`Check and apply changes to nodes via a change control system.`,
}

var changeControlCheckCmd = &cobra.Command{
	Use:	"check",
	Short:	"Review Change Control.",
	Long: `Review all unapplied changes for this node.`,
	Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(common.VERSION)
  },
}

var changeControlApplyCmd = &cobra.Command{
        Use:    "apply",
        Short:  "Apply Change Control.",
	Long: `Apply all unapplied changes to this node.`,
        Run: func(cmd *cobra.Command, args []string) {
    fmt.Println(common.VERSION)
  },
}

var registerCmd = &cobra.Command{
	Use:	"register",
	Short:	"Register node.",
	Long:	`Register this node with the Alces Monitoring service.`,
	Run:	commands.Register,
}

func main() {
	// Load config
        err := config.LoadConfig()

        if err != nil {
		return
        }

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&common.Debug, "debug", "d", false, "debug output")

	// Check if node is registered
	registered, err := coordinator.CheckRegistered()

	if err != nil {
		return
	}


	changeControlCmd.AddCommand(changeControlCheckCmd)
	changeControlCmd.AddCommand(changeControlApplyCmd)

	// Add subcommands
	if !registered {
		rootCmd.AddCommand(registerCmd)
	} else {
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
	rootCmd.Execute()
}
