package controllers

import (
	"encoding/json"
	"github.com/giovannicammarata/gc-simple-home/models"
	"log"
	"net/http"
	"runtime/debug"
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

	log.Println("command ", request.Entity, request.Target, request.Cmd)

	entityConfiguration := (*Config.Entities)[request.Entity]

	command := entityConfiguration.Commands[request.Cmd]
	targetEnv := entityConfiguration.Env[request.Target]

	var newCommand = replace(command, targetEnv)
	newCommand = replace(newCommand, Config.SystemParameters)
	log.Println("executing", newCommand)

	var code = 400
	if len(strings.TrimSpace(newCommand)) != 0 {
		_, err := http.Get(newCommand) // Send the HTTP request to the bridge

		if err == nil {
			code = 200
		} else {
			log.Println("error executing the request", body, string(debug.Stack()), err)
		}
	}

	this.Ctx.ResponseWriter.WriteHeader(code) // Return the response
}

//
// Build the command
//
func replace(str string, values map[string]string) string {
	for k, v := range values {
		str = strings.Replace(str, "{"+k+"}", v, -1)
	}

	return str
}
