package models

type SystemConfiguration struct {
	Network  NetworkConfiguration  `json:"network"`
	Domotica DomoticaConfiguration `json:"domotica"`
}

type DomoticaConfiguration struct {
	SystemParameters map[string]string              `json:"systemParameters"`
	Entities         map[string]EntityConfiguration `json:"entities"`
}

type EntityConfiguration struct {
	Env      map[string]map[string]string `json:"env"`
	Commands map[string]string            `json:"commands"`
}

type EnvConfiguration struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

type NetworkConfiguration struct {
	HTTPPort      int    `json:"http-port"`
	HTTPSPort     int    `json:"https-port"`
	HTTPSKeyFile  string `json:"https-key-file"`
	HTTPSCertFile string `json:"https-cert-file"`
	Token         string `json:"token"`
}
