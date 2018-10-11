package igcDB

type ServerInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"verion"`
}

type gilderdb struct {
	gliders map[string]Header
}
