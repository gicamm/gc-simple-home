package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/giovannicammarata/simple_home/models"
	"net/http"
	"strings"
)

type DomoticaController struct {
	BaseController
}

var Config *models.DomoticaConfiguration

//
// Executes the request by forwarding it to the bridge
//
func (this *DomoticaController) Post() {

	body := this.Ctx.Input.RequestBody
	var request models.DomoticaRequest
	json.Unmarshal(body, &request)

	fmt.Println("command ", request.Entity, request.Target, request.Cmd)

	entityConfiguration := Config.Entities[request.Entity]

	command := entityConfiguration.Commands[request.Cmd]
	targetEnv := entityConfiguration.Env[request.Target]

	var newCommand = replace(command, targetEnv)
	newCommand = replace(newCommand, Config.SystemParameters)
	fmt.Println("Executing", newCommand)

	_, err := http.Get(newCommand) // Executes the HTTP request to the bridge

	var code = 200
	if err != nil {
		code = 404
		fmt.Println(err)
	}

	this.Ctx.ResponseWriter.WriteHeader(code) // Return the response
}

//
// Build the command
//
func replace(str string, values map[string]string) string {
	for k, v := range values {
		fmt.Println(k)
		fmt.Println(v)
		str = strings.Replace(str, "{"+k+"}", v, -1)
	}

	return str
}
