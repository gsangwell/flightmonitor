package alerting

import (
	"fmt"
	"strconv"
	"flightmonitor/internal/common"
)

var service_name = "fcops-zabbix-agent2"

func ServiceIsActive() (bool, error) {
	return common.ServiceIsActive(service_name)
}

func EnableService() (bool, error) {
	return common.EnableService(service_name)
}

func DisableService() (bool, error) {
	return common.DisableService(service_name)
}

func DebugService() {
	service_active, _ := ServiceIsActive()

        fmt.Println("Service: " + service_name)
        fmt.Println("Status: " + strconv.FormatBool(service_active))
}
