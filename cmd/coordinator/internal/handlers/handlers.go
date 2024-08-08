package handlers

import (
    "net/http"
    "encoding/json"
    "strconv"
    "flightmonitor/internal/common"
    "flightmonitor/cmd/coordinator/internal/roster"
)

func Status(w http.ResponseWriter, r *http.Request) {
        response := common.StatusResponse{Version: common.VERSION, Status: "ok"}
        json.NewEncoder(w).Encode(response)
}

func Register(w http.ResponseWriter, r *http.Request) {
	server := r.URL.Query().Get("server")
	serial := r.URL.Query().Get("serial")
	ip, err := common.GetRequestIP(r)

	if err != nil {
		response := common.ResultResponse{Result: false, Error: "Failed to register server."}
                json.NewEncoder(w).Encode(response)
                return
	}

	if server != "" && serial != "" && ip != "" {
		// Sanity check node here

		// Register node
		err = roster.Register(server, serial, ip)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to register server."}
			json.NewEncoder(w).Encode(response)
			return
		} else {
			response := common.ResultResponse{Result: true}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func CheckRegistered(w http.ResponseWriter, r *http.Request) {
	server := r.URL.Query().Get("server")

	if server != "" {
		registered, err := roster.CheckRegistered(server)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to check if server is registered."}
			json.NewEncoder(w).Encode(response)
			return
		} else {
			response := common.ResultResponse{Result: registered}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
                return
	}
}

func ManagedStatus(w http.ResponseWriter, r *http.Request) {
	server := r.URL.Query().Get("server")

	if server != "" {
		status, err := roster.GetManaged(server)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Failed to get status."}
			json.NewEncoder(w).Encode(response)
		} else {
			json.NewEncoder(w).Encode(status)
		}
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
                json.NewEncoder(w).Encode(response)
	}
}

func SetManaged(w http.ResponseWriter, r *http.Request) {
	server := r.URL.Query().Get("server")
	serial := r.URL.Query().Get("serial")
	service := r.URL.Query().Get("service")
	managed := r.URL.Query().Get("managed")

	if server != "" && serial != "" && service != "" && managed != "" {
		valid_server := roster.CheckServer(server, serial)

		if !valid_server {
			response := common.ResultResponse{Result: false, Error: "Server invalid."}
			json.NewEncoder(w).Encode(response)
			return
		}

		managed_bool, err := strconv.ParseBool(managed)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
			json.NewEncoder(w).Encode(response)
			return
		}

		result_response, err := roster.SetManaged(server, service, managed_bool)

		if err != nil {
			response := common.ResultResponse{Result: false, Error: "Error setting service status."}
			json.NewEncoder(w).Encode(response)
			return
		}

		json.NewEncoder(w).Encode(result_response)
	} else {
		response := common.ResultResponse{Result: false, Error: "Invalid parameters."}
		json.NewEncoder(w).Encode(response)
		return
	}
}
