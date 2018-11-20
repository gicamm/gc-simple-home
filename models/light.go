package models

type LightRequest struct {
	Cmd    string `json:"cmd"`
	Target string `json:"target"`
}
