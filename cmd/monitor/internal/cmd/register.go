package cmd

import (
        "flightmonitor/cmd/monitor/internal/commands"
        "github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
        Use:    "register",
        Short:  "Register node.",
        Long:   `Register this node with the Alces Monitoring service.`,
        Run:    commands.Register,
}
