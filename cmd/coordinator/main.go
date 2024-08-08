package main


import (
    "net/http"
    "log/slog"
    "flightmonitor/cmd/coordinator/internal/config"
    "flightmonitor/cmd/coordinator/internal/handlers"
)

func main() {
	// Load config
	err := config.LoadConfig()

	if err != nil {
		return
	}

        // API routes
        mux := http.NewServeMux()
        mux.HandleFunc("/status", handlers.Status)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/registered", handlers.CheckRegistered)
	mux.HandleFunc("/managed/status", handlers.ManagedStatus)
	mux.HandleFunc("/managed/set", handlers.SetManaged)
        http.Handle("/", mux)

        slog.Info("Coordinator server starting on " + config.AppConfig.Coordinator.Listen + ":" + config.AppConfig.Coordinator.Port)

        http.ListenAndServe(config.AppConfig.Coordinator.Listen + ":" + config.AppConfig.Coordinator.Port, nil)
}
