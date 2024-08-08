package roster

import (
	"log/slog"
	"encoding/json"
	"strconv"
	"errors"
	"flightmonitor/internal/common"
	"flightmonitor/cmd/coordinator/internal/config"
)

func Register(server string, serial string, ip string) error {
	roster_server := config.AppConfig.Server.Protocol + "://" + config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	api_endpoint := "/register"

	site := config.AppConfig.Cluster.Site
	cluster := config.AppConfig.Cluster.Name

	url := roster_server + api_endpoint +  "?site=" + site + "&cluster=" + cluster + "&server=" + server + "&serial=" + serial + "&ip=" + ip

	response, err := common.GetApiResponse(url)

	if err != nil {
		slog.Error("roster.Register", "server", server, "serial", serial, "url", url, "err", err)
		return err
	}

	var result_response common.ResultResponse

	err = json.Unmarshal(response, &result_response)

	if err != nil {
		slog.Error("roster.Register", "server", server, "serial", serial, "url", url, "response", response, "err", err)
		return err
	}

	if result_response.Result {
		return nil
	} else {
		return errors.New(result_response.Error)
	}
}

func CheckRegistered(server string) (bool, error) {
	roster_server := config.AppConfig.Server.Protocol + "://" + config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	api_endpoint := "/registered"

	site := config.AppConfig.Cluster.Site
        cluster := config.AppConfig.Cluster.Name

        url := roster_server + api_endpoint +  "?site=" + site + "&cluster=" + cluster + "&server=" + server

	response, err := common.GetApiResponse(url)

	if err != nil {
                slog.Error("roster.CheckRegistered", "server", server, "url", url, "err", err)
                return false, err
        }

	var result_response common.ResultResponse

        err = json.Unmarshal(response, &result_response)

        if err != nil {
                slog.Error("roster.CheckRegistered", "server", server, "url", url, "response", response, "err", err)
                return false, err
        }

	return result_response.Result, nil
}

func GetManaged(server string) (*[]common.ManagedService, error) {
        roster_server := config.AppConfig.Server.Protocol + "://" + config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
        api_endpoint := "/managed/status"

        site := config.AppConfig.Cluster.Site
        cluster := config.AppConfig.Cluster.Name

        url := roster_server + api_endpoint +  "?site=" + site + "&cluster=" + cluster + "&server=" + server

        response, err := common.GetApiResponse(url)

        if err != nil {
                slog.Error("roster.GetManagedStatus", "server", server, "url", url, "err", err)
                return nil, err
        }

        var result_response []common.ManagedService

        err = json.Unmarshal(response, &result_response)

        if err != nil {
                slog.Error("roster.GetManagedStatus", "server", server, "url", url, "response", response, "err", err)
                return nil, err
        }

        return &result_response, nil
}

func SetManaged(server string, service string, managed bool) (*common.ResultResponse, error) {
	roster_server := config.AppConfig.Server.Protocol + "://" + config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
        api_endpoint := "/managed/set"

	site := config.AppConfig.Cluster.Site
        cluster := config.AppConfig.Cluster.Name

        url := roster_server + api_endpoint +  "?site=" + site + "&cluster=" + cluster + "&server=" + server + "&service=" + service + "&managed=" + strconv.FormatBool(managed)

	response, err := common.GetApiResponse(url)

	if err != nil {
                slog.Error("roster.SetManaged", "server", server, "service", service, "managed", managed, "url", url, "err", err)
                return nil, err
        }

	var result_response common.ResultResponse

	err = json.Unmarshal(response, &result_response)

	if err != nil {
                slog.Error("roster.SetManaged", "server", server, "service", service, "managed", managed, "url", url, "response", response, "err", err)
                return nil, err
        }

	return &result_response, nil
}

func GetServer(name string) (*common.Server, error) {
	roster_server := config.AppConfig.Server.Protocol + "://" + config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	api_endpoint := "/server/get"

	site := config.AppConfig.Cluster.Site
	cluster := config.AppConfig.Cluster.Name

	url := roster_server + api_endpoint +  "?site=" + site + "&cluster=" + cluster + "&name=" + name

	response, err := common.GetApiResponse(url)

	if err != nil {
		slog.Error("roster.GetServer", "name", name, "url", url, "err", err)
		return nil, err
	}

	var server common.Server

	err = json.Unmarshal(response, &server)

	if err != nil {
		slog.Error("roster.GetServer", "name", name, "url", url, "response", response, "err", err)
		return nil, err
	}

	return &server, nil
}

func CheckServer(name string, serial string) (bool) {
	server, _ := GetServer(name)

	if server != nil {
		if server.Serial == serial {
			return true
		}
	}

	return false
}
