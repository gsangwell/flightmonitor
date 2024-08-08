package database

import (
        "log/slog"
)

func (c *DBClient) AddSite(site string) error {
        _, err := c.Database.Exec("INSERT INTO sites (name) VALUES (?)", site)

        if err != nil {
                slog.Error("addSite", "site", site, "err", err)
                return err
        }

        return nil
}

func (c *DBClient) GetSiteId(site string) (int, error) {
        row := c.Database.QueryRow("SELECT id FROM sites WHERE name = ?", site)

        var id int

        err := row.Scan(&id)

        if err != nil {
                slog.Error("getSiteId", "site", site, "err", err)
                return 0, err
        }

        return id, nil
}
