package handlers

import (
	"net/http"
	"encoding/json"
	"log/slog"
	"flightmonitor/internal/common"
	"flightmonitor/cmd/roster/internal/database"
)

func AddClusterHandler(w http.ResponseWriter, r *http.Request){
	site := r.URL.Query().Get("site")
	cluster := r.URL.Query().Get("name")
	if site != "" && cluster != "" {
		err := database.Client.AddCluster(site, cluster)

		if err != nil {
                        response := common.ResultResponse{Result: false, Error: "Failed to add cluster"}
                        json.NewEncoder(w).Encode(response)
                } else {
                        response := common.ResultResponse{Result: true}
                        json.NewEncoder(w).Encode(response)
                }
	} else {
                slog.Error("addClusterHandler", "site", site, "cluster", cluster, "err", "Invalid parameters")
                response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
        }
}
