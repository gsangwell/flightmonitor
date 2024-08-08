package commands

import (
	"fmt"
	"log/slog"
	"flightmonitor/internal/common"
	"flightmonitor/cmd/monitor/internal/coordinator"
	"flightmonitor/cmd/monitor/internal/alerting"
	"flightmonitor/cmd/monitor/internal/salt"
	"github.com/spf13/cobra"
)

func EnableService(cmd *cobra.Command, args []string) {
        service := args[0]

	var is_enabled bool
	var err error

	services, err := coordinator.GetManagedStatus()

	if err != nil {
		fmt.Println("Error - unable to lookup current status")
                return
        }

        if (*services)[service].Enabled == "1" {
		fmt.Println("Service " + service + " already enabled.")
		return
	}

	switch service {
		case "alerting":
			is_enabled, err = alerting.EnableService()
		case "change-control":
			is_enabled, err = salt.EnableService()
		default:
			fmt.Println("Unknown service " + service)
			return
	}

	if err != nil || !is_enabled {
		fmt.Println("Unable to enable " + service + ".")
		if common.Debug {
			slog.Error("commands.EnableService", "service", service, "err", err)
		}
		return
	}

	is_managed, err := coordinator.SetStatus(service, true)

	if err != nil || !is_managed {
		fmt.Println("Unable to enable " + service + ".")
		if common.Debug {
                        slog.Error("commands.EnableService", "service", service, "err", err)
                }
                return
	}

	fmt.Println("Enabled " + service + " service.")
}

func DisableService(cmd *cobra.Command, args []string) {
        service := args[0]

	var is_disabled bool
	var err error

	services, err := coordinator.GetManagedStatus()

        if err != nil {
                fmt.Println("Error - unable to lookup current status")
                return
        }

        if (*services)[service].Enabled == "0" {
                fmt.Println("Service " + service + " already disabled.")
                return
        }

	switch service {
		case "alerting":
			is_disabled, err = alerting.DisableService()
		case "change-control":
			is_disabled, err = salt.DisableService()
		default:
			fmt.Println("Unknown service " + service)
			return
	}

	if err != nil || !is_disabled {
		fmt.Println("Unable to disable " + service + ".")
		return
	}

	is_unmanaged, err := coordinator.SetStatus(service, false)

	if err != nil || !is_unmanaged {
		fmt.Println("Unable to disable " + service + ".")
		return
	}

	fmt.Println("Disabled " + service + " service.")
}

func StatusService(cmd *cobra.Command, args []string) {
	service := args[0]

	is_managed, err := coordinator.GetStatus(service)

	if err != nil {
		fmt.Println("Unable to check status of " + service)
	} else {
		if is_managed {
			fmt.Println(service + ": enabled")
		} else {
			fmt.Println(service + ": disabled")
		}
	}
}
