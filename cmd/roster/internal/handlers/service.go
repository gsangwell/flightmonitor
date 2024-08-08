package handlers

import (
        "net/http"
        "encoding/json"
	"strconv"
	"log/slog"
	"errors"
        "flightmonitor/internal/common"
	"flightmonitor/cmd/roster/internal/database"
	"flightmonitor/cmd/roster/internal/slack"
	"flightmonitor/cmd/roster/internal/zabbix"
)

func GetManagedStatusHandler(w http.ResponseWriter, r *http.Request){
	site := r.URL.Query().Get("site")
	cluster := r.URL.Query().Get("cluster")
	server := r.URL.Query().Get("server")
	service := r.URL.Query().Get("service")

	if site != "" && cluster != "" && server != "" && service != "" {
		managed_status, err := database.Client.GetManagedStatus(site, cluster, server, service)

		if managed_status != nil {
			json.NewEncoder(w).Encode(managed_status)
		} else {
			response := common.ResultResponse{Result: false, Error: err.Error()}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
		json.NewEncoder(w).Encode(response)
	}
}

func SetManagedStatusHandler(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
        cluster := r.URL.Query().Get("cluster")
        server := r.URL.Query().Get("server")
        service := r.URL.Query().Get("service")
	managed_string := r.URL.Query().Get("managed")
	managed, err := strconv.ParseBool(managed_string)

        if site != "" && cluster != "" && server != "" && service != "" && err == nil {
		var err error
		var hostid string

		switch service {
			case "alerting":
				hostid, err = zabbix.Client.GetHostId(site, cluster, server)

				if err != nil {
					break
				}
				err = zabbix.Client.SetEnable(hostid, managed)
	                case "change-control":
				err = nil
			default:
				err = errors.New("Unknown service " + service + ".")
		}

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to enable " + service + "."}
			json.NewEncoder(w).Encode(response)
			return
		}

		err = database.Client.SetManagedStatus(site, cluster, server, service, managed)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to update managed status"}
			json.NewEncoder(w).Encode(response)
		} else {
			slack.Client.SendMessage("Site: " + site + "\n" + "Cluster: " + cluster + "\n" + "Host: " + server + "\n" + "Service: " + service + "\n" + "Enabled: " + managed_string)
			response := common.ResultResponse{Result: true}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		slog.Error("setManagedStatusHandler", "site", site, "cluster", cluster, "server", server, "service", service, "managed", managed, "err", err)
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
	}
}

func GetAllManagedServices(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
        cluster := r.URL.Query().Get("cluster")
        server := r.URL.Query().Get("server")

	if site != "" && cluster != "" && server != "" {
		services, err := database.Client.GetAllManagedServices(site, cluster, server)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: err.Error()}
                        json.NewEncoder(w).Encode(response)
		} else {
			json.NewEncoder(w).Encode(services)
		}
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
	}
}
