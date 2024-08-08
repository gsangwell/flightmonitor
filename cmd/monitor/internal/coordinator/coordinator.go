package coordinator

import (
	"errors"
	"encoding/json"
	"strconv"
	"log/slog"
	"flightmonitor/internal/common"
	"flightmonitor/cmd/monitor/internal/config"
	"flightmonitor/cmd/monitor/internal/system"
)

func CheckRegistered() (bool, error) {
	hostname, _ := system.GetHostname()

	if hostname == nil {
                return false, errors.New("Unable to determine hostname.")
        }

	coordinator_server := config.CliConfig.Server.Protocol + "://" + config.CliConfig.Server.Host + ":" + config.CliConfig.Server.Port
	api_endpoint := "/registered"
	url := coordinator_server + api_endpoint + "?server=" + *hostname

	response, err := common.GetApiResponse(url)

	if err != nil {
                if common.Debug {
                        slog.Error("coordinator.CheckRegistered", "url", url, "err", err)
                }
                return false, errors.New("Unable to check registered status.")
        }

        var result_response common.ResultResponse

        err = json.Unmarshal(response, &result_response)

        if err != nil {
                if common.Debug {
                        slog.Error("coordinator.CheckRegistered", "url", url, "response", response, "err", err)
		}
                return false, errors.New("Unable to check registered status.")
        }

	return result_response.Result, nil
}

func Register() (bool, error) {
	hostname, _ := system.GetHostname()
	serial, _ := system.GetSerial()

        if hostname == nil {
                return false, errors.New("Unable to determine hostname.")
        }

        if serial == nil {
                return false, errors.New("Unable to determine serial.")
        }

	coordinator_server := config.CliConfig.Server.Protocol + "://" + config.CliConfig.Server.Host + ":" + config.CliConfig.Server.Port
        api_endpoint := "/register"
	url := coordinator_server + api_endpoint + "?server=" + *hostname + "&serial=" + *serial

	response, err := common.GetApiResponse(url)

	if err != nil {
                if common.Debug {
                        slog.Error("coordinator.Register", "url", url, "err", err)
		}
                return false, errors.New("Unable to register server.")
        }

	var result_response common.ResultResponse
        err = json.Unmarshal(response, &result_response)

	if err != nil {
                if common.Debug {
                        slog.Error("coordinator.Register", "url", url, "response", response, "err", err)
                }
                return false, errors.New("Unable to register server.")
        }

	return result_response.Result, nil
}

func GetStatus(service string) (bool, error) {
	services, err := GetManagedStatus()

	if err != nil {
		return false, errors.New("Unable to get services.")
	}

	svc, ok := (*services)[service]

	if ok {
		enabled, err := strconv.ParseBool(svc.Enabled)

		if err != nil {
			return false, errors.New("Unable to get status of service.")
		}

		return enabled, nil
	} else {
		return false, errors.New("Unable to find service.")
	}
}

func SetStatus(service string, enabled bool) (bool, error) {
        hostname, _ := system.GetHostname()
        serial, _ := system.GetSerial()

	if hostname == nil {
                return false, errors.New("Unable to determine hostname.")
        }

        if serial == nil {
                return false, errors.New("Unable to determine serial.")
        }

	coordinator_server := config.CliConfig.Server.Protocol + "://" + config.CliConfig.Server.Host + ":" + config.CliConfig.Server.Port
        api_endpoint := "/managed/set"
        url := coordinator_server + api_endpoint + "?server=" + *hostname + "&serial=" + *serial + "&service=" + service + "&managed=" + strconv.FormatBool(enabled)

        response, err := common.GetApiResponse(url)

	if err != nil {
		if common.Debug {
                        slog.Error("coordinator.SetStatus", "service", service, "err", err)
                }
                return false, errors.New("Unable to set current status.")
        }

        var result_response common.ResultResponse

        err = json.Unmarshal(response, &result_response)

	if err != nil {
		if common.Debug {
                        slog.Error("coordinator.SetStatus", "service", service, "enabled", enabled, "response", response, "err", err)
                }
                return false, errors.New("Unable to set current status.")
        }

	if !result_response.Result {
		if common.Debug {
                        slog.Error("coordinator.SetStatus", "service", service, "enabled", enabled, "response", response, "err", err)
                }
		return false, errors.New("Unable to set current status.")
	}

	return true, nil
}

func GetManagedStatus() (*map[string]*common.ManagedService, error) {
	hostname, _ := system.GetHostname()
        serial, _ := system.GetSerial()

        if hostname == nil {
                return nil, errors.New("Unable to determine hostname.")
        }

        if serial == nil {
                return nil, errors.New("Unable to determine serial.")
        }

        coordinator_server := config.CliConfig.Server.Protocol + "://" + config.CliConfig.Server.Host + ":" + config.CliConfig.Server.Port
        api_endpoint := "/managed/status"
        url := coordinator_server + api_endpoint + "?server=" + *hostname

        response, err := common.GetApiResponse(url)

        if err != nil {
                if common.Debug {
                        slog.Error("coordinator.GetManagedStatus", "err", err)
                }
                return nil, errors.New("Unable to get current status.")
        }

        var managed_status []common.ManagedService

        err = json.Unmarshal(response, &managed_status)

        if err != nil {
                if common.Debug {
                        slog.Error("coordinator.GetManagedStatus", "response", response, "err", err)
                }
                return nil, errors.New("Unable to get current status.")
        }

	managed_services := make(map[string]*common.ManagedService)

	for _, service := range managed_status {
		svc := common.ManagedService{Service: service.Service, Enabled: service.Enabled}
		managed_services[service.Service] = &svc
	}

        return &managed_services, nil
}
