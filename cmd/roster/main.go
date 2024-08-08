package main


import (

    "log/slog"
    "net/http"

    "flightmonitor/cmd/roster/internal/config"
    "flightmonitor/cmd/roster/internal/slack"
    "flightmonitor/cmd/roster/internal/zabbix"
    "flightmonitor/cmd/roster/internal/database"
    "flightmonitor/cmd/roster/internal/handlers"

    _ "github.com/go-sql-driver/mysql"

)

func main() {
	// Config
	err := config.LoadConfig()

        if err != nil {
                return
        }

	// Database
	err = database.Init(config.RosterConfig.Database.Host, config.RosterConfig.Database.Port, config.RosterConfig.Database.Username, config.RosterConfig.Database.Password, config.RosterConfig.Database.Database)

	if err != nil {
		return
	}

	// Zabbix
	err = zabbix.Init(config.RosterConfig.Zabbix.API, config.RosterConfig.Zabbix.Username, config.RosterConfig.Zabbix.Password)

	if err != nil {
		return
	}

	// Slack
	err = slack.Init(config.RosterConfig.Slack.Token, config.RosterConfig.Slack.ChannelID)

	if err != nil {
		return
	}

	// Static files
	//fs := http.FileServer(http.Dir("build"))
	//http.Handle("/", fs)

	// API routes
	api := http.NewServeMux()
	api.HandleFunc("/status", handlers.StatusHandler)

	api.HandleFunc("/register", handlers.Register)
	api.HandleFunc("/registered", handlers.CheckRegistered)

	api.HandleFunc("/site/add", handlers.AddSiteHandler)


	api.HandleFunc("/cluster/add", handlers.AddClusterHandler)

	api.HandleFunc("/server/get", handlers.GetServer)

	api.HandleFunc("/managed/status", handlers.GetAllManagedServices)
	api.HandleFunc("/managed/set", handlers.SetManagedStatusHandler)
	api.HandleFunc("/managed/get", handlers.GetManagedStatusHandler)

	http.Handle("/", api)

	slog.Info("Roster server starting on " + config.RosterConfig.Roster.Listen + ":" + config.RosterConfig.Roster.Port)

	http.ListenAndServe(config.RosterConfig.Roster.Listen + ":" + config.RosterConfig.Roster.Port, nil)
}
