package database

import (
        "log/slog"
	"github.com/jmoiron/sqlx"
)

type DBClient struct {
	Database *sqlx.DB
}

var Client *DBClient

func Init(host string, port string, username string, password string, database string) error {
	db, err := sqlx.Open("mysql", username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database)

        if err != nil {
                slog.Error(err.Error())
                return err
        }

        if err := db.Ping(); err != nil {
                slog.Error(err.Error())
                return err
        }

        slog.Info("Connected to database.")

	Client = &DBClient {
		Database: db,
	}

	return nil
}
