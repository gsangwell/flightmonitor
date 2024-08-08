package common

type StatusResponse struct {
        Version string `json:"version"`
        Status string `json:"status"`
}

type ResultResponse struct {
	Result bool `json:"result"`
	Error string `json:"error"`
}
