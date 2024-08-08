package common

import (
	"net/http"
	"encoding/json"
	"time"
	"io"
	"errors"
	"log/slog"
	"strings"
	"net"
)

func GetApiResponse(url string) ([]byte, error) {
        httpClient := http.Client {
                Timeout: time.Second * 2,
        }

        req, err := http.NewRequest(http.MethodGet, url, nil)

        if err != nil {
		if Debug {
                        slog.Error("common.GetApiResponse", "url", url, "err", err)
                }
		return nil, errors.New("Failed to create API request.")
        }

        res, err := httpClient.Do(req)

        if err != nil {
		if Debug {
                        slog.Error("common.GetApiResponse", "url", url, "res", res, "err", err)
                }
		return nil, errors.New("Failed to call API request.")
        }

        if res.Body != nil {
                defer res.Body.Close()
        }

        body, err := io.ReadAll(res.Body)

        if err != nil {
		if Debug {
                        slog.Error("common.GetApiResponse", "url", url, "res", res, "body", body, "err", err)
                }
		return nil, errors.New("Failed to retrieve API response.")
        }

        var result_response ResultResponse
        err = json.Unmarshal(body, &result_response)

	if err == nil && result_response.Error != "" {
		slog.Error("common.GetApiResponse", "url", url, "res", res, "body", body, "result", result_response, "err", err)
		return nil, errors.New(result_response.Error)
	}

	return body, nil
}

func GetRequestIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// get last IP in list since ELB prepends other user defined IPs, meaning the last one is the actual client IP.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
