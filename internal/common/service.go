package common

import (
	"strings"
	"errors"
)

func ServiceIs(service_name string, status string) (bool, error) {
	exit_code, stdout, _, _ := RunCommand("systemctl", "show", service_name, "-p", "ActiveState")

        if exit_code != 0 {
                return false, errors.New("Failed to check systemd service.")
        }

        stdout = strings.Replace(stdout, "\n", "", -1)
	current_status := strings.Split(stdout, "=")[1]

        if current_status == status {
                return true, nil
        } else {
                return false, nil
        }
}

func ServiceIsActive(service_name string) (bool, error) {
	return ServiceIs(service_name, "active")
}

func ServiceIsInactive(service_name string) (bool, error) {
	return ServiceIs(service_name, "inactive")
}

func EnableService(service_name string) (bool, error) {
	exit_code, _, _, _ := RunCommand("systemctl", "enable", "--now", service_name)

	if exit_code != 0 {
		return false, errors.New("Failed to enable systemd service.")
	}

	return ServiceIsActive(service_name)
}

func DisableService(service_name string) (bool, error) {
	exit_code, _, _, _ := RunCommand("systemctl", "disable", "--now", service_name)

	if exit_code != 0 {
		return false, errors.New("Failed to disable systemd service.")
	}

	return ServiceIsInactive(service_name)
}
