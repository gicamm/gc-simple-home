package models

type AuthorizedRequest struct {
	Token string `json:"token"`
}

type DomoticaRequest struct {
	AuthorizedRequest
	Entity string `json:"entity"`
	Cmd    string `json:"cmd"`
	Target string `json:"target"`
}
