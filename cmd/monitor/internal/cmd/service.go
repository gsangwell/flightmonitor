package cmd

import (
	"strings"
        "flightmonitor/cmd/monitor/internal/commands"
        "github.com/spf13/cobra"
)

var valid_services = []string{"alerting", "change-control"}

var statusCmd = &cobra.Command{
        Use:    "status <" + strings.Join(valid_services, "|") + ">",
        Short:  "Service Status",
        Long:   `Show whether a particular monitoring service is enabled or disabled on this host.`,
        Args:   cobra.ExactValidArgs(1),
        ValidArgs: valid_services,
        Run:    commands.StatusService,
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
