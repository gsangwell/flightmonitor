package salt

import (
	"flightmonitor/internal/common"
)

var service_name = "fcops-salt-minion"

func ServiceIsActive() (bool, error) {
	return common.ServiceIsActive(service_name)
}

func EnableService() (bool, error) {
	return common.EnableService(service_name)
}

func DisableService() (bool, error) {
	return common.DisableService(service_name)
}
