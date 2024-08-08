package common

type Site struct {
        Id int `json:"id"`
        Name string `json:"name"`
}

type Cluster struct {
        Id int `json:"id"`
        Site string `json:"site"`
        Name string `json:"name"`
}

type Server struct {
        Id int `json:"id"`
        Site string `json"site"`
        Cluster string `json"cluster"`
        Name string `json"name"`
	Serial string `json"serial"`
	IP string `json"ip"`
}

type Service struct {
        Id int `json:"id"`
        Name string `json"name"`
}

type ManagedService struct {
	Service string `json:"service"`
	Enabled string `json:"enabled"`
}

type ManagedStatus struct {
        Site string `json:"site"`
        Cluster string `json:"cluster"`
        Node string `json:"node"`
        Service string `json:"service"`
        Managed bool `json:"managed"`
}
