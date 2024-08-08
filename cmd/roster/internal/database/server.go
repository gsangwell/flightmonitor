package database

import (
        "log/slog"
	"flightmonitor/internal/common"
	"database/sql"
)

func (c *DBClient) AddServer(site string, cluster string, server string, serial string, ip string) error {
        site_id, err := c.GetSiteId(site)

        if err != nil {
                slog.Error("addServer", "site", site, "cluster", cluster, "server", server, "serial", serial, "ip", ip, "err", err)
                return err
        }

        cluster_id, err := c.GetClusterId(site, cluster)

        if err != nil {
                slog.Error("addServer", "site", site, "cluster", cluster, "server", server, "serial", serial, "ip", ip, "err", err)
                return err
        }

        services, err := c.ListServices()

        if err != nil {
                slog.Error("addServer", "site", site, "cluster", cluster, "server", server, "serial", serial, "ip", ip, "err", err)
                return err
        }

        row := c.Database.QueryRow("INSERT INTO servers (site, cluster, name, serial, ip) VALUES(?, ?, ?, ?, ?) RETURNING id", site_id, cluster_id, server, serial, ip)

        var server_id int
        err = row.Scan(&server_id)

        if err != nil {
                slog.Error("addServer", "site", site, "cluster", cluster, "server", server, "serial", serial, "ip", ip, "err", err)
                return err
        }

        //server_id := result.LastInsertId

        for _,service := range services {
                _, err = c.Database.Exec("INSERT INTO management_status (server, service, managed) VALUES (?, ?, false)", server_id, service.Id)

                if err != nil {
                        slog.Error("addServer", "site", site, "cluster", cluster, "server", server, "serial", serial, "ip", ip, "err", err)
                }
        }

        return nil
}

func (c *DBClient) ServerExists(site string, cluster string, server string) (bool, error) {
	row := c.Database.QueryRow("SELECT COUNT(*) FROM servers JOIN clusters ON clusters.id = servers.cluster JOIN sites ON sites.id = clusters.site WHERE sites.name = ? AND clusters.name = ? AND servers.name = ?", site, cluster, server)

	var count int
	err := row.Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		slog.Error("database.ServerExists", "site", site, "cluster", cluster, "server", server, "err", err)
		return false, err
	} else {
		return count != 0, nil
	}
}

func (c *DBClient) GetServerId(site string, cluster string, server string) (int, error) {
        row := c.Database.QueryRow("SELECT servers.id FROM servers JOIN clusters ON clusters.id = servers.cluster JOIN sites ON sites.id = clusters.site WHERE sites.name = ? AND clusters.name = ? AND servers.name = ?", site, cluster, server)

        var id int

        err := row.Scan(&id)

        if err != nil {
                slog.Error("getServerId", "site", site, "cluster", cluster, "server", server, "err", err)
                return 0, err
        }

        return id, nil
}

func (c *DBClient) GetServer(site string, cluster string, name string) (*common.Server, error) {
	row := c.Database.QueryRowx("SELECT * FROM servers JOIN clusters ON clusters.id = servers.cluster JOIN sites ON sites.id = clusters.site WHERE sites.name = ? AND clusters.name = ? AND servers.name = ?", site, cluster, name)

	var server common.Server

	err := row.StructScan(&server)

	if err != nil {
		slog.Error("database.getServer", "site", site, "cluster", cluster, "server", name, "err", err)
		return nil, err
	}

	return &server, nil
}
