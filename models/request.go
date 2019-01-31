package models

type DomoticaRequest struct {
	Entity string `json:"entity"`
	Cmd    string `json:"cmd"`
	Target string `json:"target"`
}
