package handlers

import (
        "net/http"
        "encoding/json"
	"log/slog"
        "flightmonitor/internal/common"
	"flightmonitor/cmd/roster/internal/database"
)

func AddSiteHandler(w http.ResponseWriter, r *http.Request){
	site := r.URL.Query().Get("name")

	if site != "" {
		err := database.Client.AddSite(site)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to add site"}
			json.NewEncoder(w).Encode(response)
		} else {
			response := common.ResultResponse{Result: true}
			json.NewEncoder(w).Encode(response)
		}
	} else {
                slog.Error("addSiteHandler", "site", site, "err", "Invalid parameters")
                response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
        }
}
