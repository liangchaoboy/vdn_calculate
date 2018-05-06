package server

type Domains struct {
	Items []DomainInfo `json:"items"`
}

type DomainInfo struct {
	Name    string       `json:"name"`
	Domains []DomainItem `json:"domains`
}

type DomainItem struct {
	Domain string `json:"domain`
	Type   string `json:"type"`
}
