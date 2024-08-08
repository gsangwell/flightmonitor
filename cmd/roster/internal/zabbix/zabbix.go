package zabbix

import (
	"log/slog"
	"errors"

	gozabbix "github.com/claranet/go-zabbix-api"
)

type ZabbixClient struct {
        ZabbixClient *gozabbix.API
}

var Client *ZabbixClient

func Init(api_url string, username string, password string) error {

	zc, err := gozabbix.NewAPI(api_url)

	if err != nil {
		slog.Error("zabbix.Init", "api_url", api_url, "username", username, "password", password, "err", err)
		return err
	}

	_, err = zc.Login(username, password)

	if err != nil {
		slog.Error("zabbix.Init", "api_url", api_url, "username", username, "password", password, "err", err)
		return err
	}

	res, err := zc.Version()

	if err != nil {
		slog.Error("zabbix.Init", "api_url", api_url, "username", username, "password", password, "err", err)
		return err
	}

	Client = &ZabbixClient {
                ZabbixClient: zc,
        }

	slog.Info("zabbix.Init", "api_version", res)

	return nil
}

func (z *ZabbixClient) SetEnable(host_id string, enabled bool) error {
	var params gozabbix.Params

	if enabled {
		params = gozabbix.Params{"hostid": host_id, "status": 0}
	} else {
		params = gozabbix.Params{"hostid": host_id, "status": 1}
	}

	res, err := z.ZabbixClient.Call("host.update", params)

	if err != nil {
                slog.Error("zabbix.SetEnable", "host_id", host_id, "enabled", enabled, "err", err)
                return err
        }

	if res.Error != nil {
		slog.Error("zabbix.SetEnable", "host_id", host_id, "enabled", enabled, "res", res)
		return err
	}

        return nil
}

func (z *ZabbixClient) GetGroupId(site string, cluster string) (string, error) {
	params := gozabbix.Params{"output": "extend", "filter": map[string]string{"name": site + "." + cluster}, "limit": 1}

	res, err := z.ZabbixClient.Call("hostgroup.get", params)

	if err != nil {
		slog.Error("zabbix.GetGroupId", "site", site, "cluster", cluster, "err", err)
		return "", err
	}

	f := res.Result.([]interface{})
	if len(f) <= 0 {
		slog.Error("zabbix.GetGroupId", "site", site, "cluster", cluster, "err", "Failed to get group ID.")
		return "", errors.New("Failed to get group ID.")
	}
	d := f[0].(map[string]interface{})
	str, ok := d["groupid"].(string)

	if !ok {
		slog.Error("zabbix.GetGroupId", "site", site, "cluster", cluster, "err", "Failed to get group ID.")
		return "", errors.New("Failed to get group ID.")
	}

	return str, nil
}

func (z *ZabbixClient) GetProxyId(site string, cluster string) (string, error) {
        params := gozabbix.Params{"output": "extend", "filter": map[string]string{"name": "proxy." + site + "." + cluster}, "limit": 1}

        res, err := z.ZabbixClient.Call("proxy.get", params)

        if err != nil {
                slog.Error("zabbix.GetProxyId", "site", site, "cluster", cluster, "err", err)
                return "", err
        }

        f := res.Result.([]interface{})
        if len(f) <= 0 {
                slog.Error("zabbix.GetProxyId", "site", site, "cluster", cluster, "err", "Failed to get proxy ID.")
                return "", errors.New("Failed to get proxy ID.")
        }
        d := f[0].(map[string]interface{})
        str, ok := d["proxyid"].(string)

        if !ok {
                slog.Error("zabbix.GetProxyId", "site", site, "cluster", cluster, "err", "Failed to get proxy ID.")
                return "", errors.New("Failed to get proxy ID.")
        }

        return str, nil
}

func (z *ZabbixClient) GetHostId(site string, cluster string, host string) (string, error) {
	g, err := z.GetGroupId(site, cluster)

        if err != nil {
                return "", err
        }

        //groups := gozabbix.HostGroupIDs{gozabbix.HostGroupID{GroupID: g}}

	params := gozabbix.Params{"output": "extend", "groupids": g, "filter": map[string]string{"name": host}, "limit": 1}

        res, err := z.ZabbixClient.Call("host.get", params)

        if err != nil {
                slog.Error("zabbix.GetHostId", "site", site, "cluster", cluster, "host", host, "err", err, "params", params)
                return "", err
        }

        f := res.Result.([]interface{})
        if len(f) <= 0 {
                slog.Error("zabbix.GetHostId", "site", site, "cluster", cluster, "host", host, "err", "Failed to get host ID.", "params", params, "res", res)
                return "", errors.New("Failed to get host ID.")
        }
        d := f[0].(map[string]interface{})
        str, ok := d["hostid"].(string)

        if !ok {
		slog.Error("zabbix.GetHostId", "site", site, "cluster", cluster, "host", host, "err", "Failed to get host ID.", "params", params, "res", res)
                return "", errors.New("Failed to get host ID.")
        }

        return str, nil
}

func (z *ZabbixClient) AddHost(site string, cluster string, server string, ip string) error {
	i := gozabbix.HostInterface{DNS: "", IP: ip, Main: 1, Port: "10050", Type: 1, UseIP: 1}
	interfaces := gozabbix.HostInterfaces{i}

	g, err := z.GetGroupId(site, cluster)

	if err != nil {
		return errors.New("Failed to add host.")
	}

	groups := gozabbix.HostGroupIDs{gozabbix.HostGroupID{GroupID: g}}

	p, err := z.GetProxyId(site, cluster)

	if err != nil {
		return errors.New("Failed to add host.")
	}

	t := gozabbix.TemplateID{TemplateID: "10001"}
	templates := gozabbix.TemplateIDs{t}

	params := gozabbix.Params{"host": server, "status": 1, "proxyid": p, "monitored_by": 1, "interfaces": interfaces, "groups": groups, "templates": templates}

	res, err := z.ZabbixClient.Call("host.create", params)

	if err != nil {
		slog.Error("zabbix.AddHost", "site", site, "cluster", cluster, "server", server, "ip", ip, "params", params, "err", err)
		return errors.New("Failed to add host.")
	}

	slog.Info("zabbix.AddHost", "res", res)

	if res.Error != nil {
		slog.Error("zabbix.AddHost", "site", site, "cluster", cluster, "server", server, "ip", ip, "params", params, "res", res)
		return errors.New("Failed to add host.")
	}

	return nil
}
