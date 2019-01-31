package models

type DomoticaConfiguration struct {
	SystemParameters map[string]string              `json:"systemParameters"`
	Entities         map[string]EntityConfiguration `json:"entities"`
}

type EntityConfiguration struct {
	//Env map[string]EnvConfiguration						`json:"env"`
	Env map[string]map[string]string `json:"env"`
	//Env map[string]interface{}								`json:"env"`
	Commands map[string]string `json:"commands"`
}

type EnvConfiguration struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}
