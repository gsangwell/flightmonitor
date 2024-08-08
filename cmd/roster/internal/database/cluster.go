package database

import (
	"log/slog"
)

func (c *DBClient) AddCluster(site string, cluster string) error {
        site_id, err := c.GetSiteId(site)

        if err != nil {
                slog.Error("addCluster", "site", site, "cluster", cluster, "err", err)
                return err
        }

        _, err = c.Database.Exec("INSERT INTO clusters (site, name) VALUES(?, ?)", site_id, cluster)

        if err != nil {
                slog.Error("addCluster", "site", site, "cluster", cluster, "err", err)
                return err
        }

        return nil
}

func (c *DBClient) GetClusterId(site string, cluster string) (int, error) {
        row := c.Database.QueryRow("SELECT clusters.id FROM clusters JOIN sites ON sites.id = clusters.site WHERE sites.name = ? AND clusters.name = ?", site, cluster)

        var id int

        err := row.Scan(&id)

        if err != nil {
                slog.Error("getClusterId", "site", site, "cluster", cluster, "err", err)
                return 0, err
        }

        return id, nil
}
