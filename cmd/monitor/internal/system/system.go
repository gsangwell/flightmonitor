package system

import (
	"strings"
	"errors"
	"os"
	"flightmonitor/internal/common"
)

func GetSerial() (*string, error) {
	exitcode, stdout, _, err := common.RunCommand("dmidecode", "-s", "system-serial-number")

	if exitcode != 0 {
		return nil, err
	}

	serial := strings.Replace(stdout, "\n", "", -1)

	return &serial, nil
}

func GetHostname() (*string, error) {
	host, err := os.Hostname()
        hostname := strings.Split(host, ".")[0]

        if err != nil || hostname == "" {
                return nil, errors.New("Unable to determine hostname")
        }

	return &hostname, nil
}
