package database

import (
        "log/slog"
	"errors"
        "flightmonitor/internal/common"
)

func (c *DBClient) ListServices() ([]common.Service, error) {
        var services []common.Service

        rows, err := c.Database.Query("SELECT services.id, services.name FROM services")

        if err != nil {
                slog.Error("listServices", "err", err)
                return nil, err
        }

        defer rows.Close()

        for rows.Next(){
                var service common.Service
                err := rows.Scan(&service.Id, &service.Name)

                if err != nil {
                        slog.Error("listServices", "err", err)
                        return nil, err
                }

                services = append(services, service)
        }

        return services, nil
}

func (c *DBClient) GetManagedId(site string, cluster string, server string, service string) (int, error) {
        row := c.Database.QueryRow("SELECT management_status.id FROM management_status JOIN servers on management_status.server = servers.id JOIN services on management_status.service = services.id JOIN clusters ON clusters.id = servers.cluster JOIN sites ON sites.id = clusters.site WHERE sites.name = ? AND clusters.name = ? AND servers.name = ? AND services.name = ? LIMIT 1;", site, cluster, server, service)

        var id int

        err := row.Scan(&id)

        if err != nil {
                slog.Error("getManagedId", "site", site, "cluster", cluster, "server", server, "service", service, "err", err)
                return 0, err
        }

        return id, nil
}

func (c *DBClient) GetManagedStatus(site string, cluster string, server string, service string) (*common.ManagedStatus, error) {
        row := c.Database.QueryRow("SELECT sites.name as site, clusters.name as cluster, servers.name as node, services.name as service, management_status.managed FROM management_status JOIN servers on management_status.server = servers.id JOIN services on management_status.service = services.id JOIN clusters ON clusters.id = servers.cluster JOIN sites ON sites.id = clusters.site WHERE sites.name = ? AND clusters.name = ? AND servers.name = ? AND services.name = ? LIMIT 1;", site, cluster, server, service)

        var managed_status common.ManagedStatus
        err := row.Scan(&managed_status.Site, &managed_status.Cluster, &managed_status.Node, &managed_status.Service, &managed_status.Managed)

        if err != nil {
                slog.Error("getManagedStatus", "site", site, "cluster", cluster, "server", server, "service", service, "err", err)
		e := errors.New("Invalid parameters.")
                return nil, e
        }

        return &managed_status, nil
}

func (c *DBClient) SetManagedStatus(site string, cluster string, server string, service string, managed bool) (error) {
        id, err := c.GetManagedId(site, cluster, server, service)

        if err != nil {
                return err
        }

        _, err = c.Database.Exec("UPDATE management_status SET managed = ? WHERE id = ?", managed, id)

        if err != nil {
                slog.Error("setManagedStatus", "site", site, "cluster", cluster, "server", server, "service", service, "err", err)
                return err
        }

        return nil
}

func (c *DBClient) GetAllManagedServices(site string, cluster string, server string) (*[]common.ManagedService, error) {
	server_id, err := c.GetServerId(site, cluster, server)

	if err != nil {
		return nil, err
	}

	rows, err := c.Database.Queryx("SELECT services.name AS service, management_status.managed AS enabled FROM management_status JOIN services ON management_status.service = services.id WHERE server = ?", server_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var services []common.ManagedService

	for rows.Next() {
		var service common.ManagedService

		err := rows.StructScan(&service)

		if err != nil {
			return &services, err
		}

		services = append(services, service)
	}

	err = rows.Err()

	if err != nil {
		return &services, err
	}

	return &services, nil
}
