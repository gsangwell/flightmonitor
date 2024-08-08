package cmd

import (
	"fmt"
	"flightmonitor/internal/common"
        "github.com/spf13/cobra"
)

var changeControlCmd = &cobra.Command{
        Use:    "change-control",
        Short:  "Change Control",
        Long:   `Check and apply changes to nodes via a change control system.`,
}

var changeControlCheckCmd = &cobra.Command{
        Use:    "check",
        Short:  "Review Change Control.",
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

func init() {
	changeControlCmd.AddCommand(changeControlCheckCmd)
        changeControlCmd.AddCommand(changeControlApplyCmd)
}
