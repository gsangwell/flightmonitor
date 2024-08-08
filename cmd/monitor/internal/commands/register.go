package commands

import (
        "fmt"
        "flightmonitor/cmd/monitor/internal/coordinator"
	"github.com/spf13/cobra"
)

func Register(cmd *cobra.Command, args []string) {
	registered, err := coordinator.Register()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if registered {
		fmt.Println("Registered node!")
		return
	} else {
		fmt.Println("Failed to register node.")
	}
}
