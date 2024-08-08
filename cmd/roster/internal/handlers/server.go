package handlers

import (
        "net/http"
        "encoding/json"
	"log/slog"
        "flightmonitor/internal/common"
	"flightmonitor/cmd/roster/internal/database"
	"flightmonitor/cmd/roster/internal/slack"
	"flightmonitor/cmd/roster/internal/zabbix"
)

func Register(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
        cluster := r.URL.Query().Get("cluster")
        server := r.URL.Query().Get("server")
        serial := r.URL.Query().Get("serial")
	ip := r.URL.Query().Get("ip")

	if site != "" && cluster != "" && server != "" && serial != "" && ip != "" {
		err := zabbix.Client.AddHost(site, cluster, server, ip)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to add server"}
			json.NewEncoder(w).Encode(response)
                        return
		}

                err = database.Client.AddServer(site, cluster, server, serial, ip)

                if err != nil {
                        response := common.ResultResponse{Result: false, Error: "Failed to add server"}
                        json.NewEncoder(w).Encode(response)
			return
                }

		slack.Client.SendMessage("New Server Registered! \nSite: " + site + "\n" + "Cluster: " + cluster + "\n" + "Host: " + server + "\nSerial:" + serial + "\nIP: " + ip)
                response := common.ResultResponse{Result: true}
                json.NewEncoder(w).Encode(response)
		return
        } else {
                slog.Error("addServerHandler", "site", site, "cluster", cluster, "server", server, "serial", serial, "err", "Invalid parameters")
                response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
        }
}

func CheckRegistered(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
        cluster := r.URL.Query().Get("cluster")
        server := r.URL.Query().Get("server")

	if site != "" && cluster != "" && server != "" {
		registered, err := database.Client.ServerExists(site, cluster, server)

		if err != nil {
			slog.Error("handlers.CheckRegistered", "site", site, "cluster", cluster, "server", server, "err", err)
			response := common.ResultResponse{Result: false, Error: "Failed to add server"}
                        json.NewEncoder(w).Encode(response)
			return
		} else {
			response := common.ResultResponse{Result: registered}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		slog.Error("handlers.CheckRegistered", "site", site, "cluster", cluster, "server", server, "err", "Invalid parameters")
                response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
		return
	}
}

func GetServer(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	cluster := r.URL.Query().Get("cluster")
	name := r.URL.Query().Get("name")

	if site != "" && cluster != "" && name != "" {
		server, _ := database.Client.GetServer(site, cluster, name)

		if server != nil {
			json.NewEncoder(w).Encode(server)
		} else {
			response := common.ResultResponse{Result: false, Error: "Failed to get server."}
			json.NewEncoder(w).Encode(response)
		}
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
		json.NewEncoder(w).Encode(response)
	}
}
