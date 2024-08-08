package cmd

import (
        "flightmonitor/cmd/coordinator/internal/config"
        "github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
        Use:    "configure",
        Short:  "(Re)Configure Coordinator.",
        Long:   `(Re)Configure the settings for Alces Monitoring Coordinator API`,
        Run:    config.Configure,
}
