package handlers

import (
        "net/http"
        "encoding/json"
        "flightmonitor/internal/common"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
        response := common.StatusResponse{Version: "v0.0.1", Status: "ok"}
        json.NewEncoder(w).Encode(response)
}
