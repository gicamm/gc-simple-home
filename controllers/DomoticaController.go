package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/giovannicammarata/simple_home/models"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type DomoticaController struct {
	BaseController
}

var config models.DomoticaConfiguration

func init() {
	jsonFile, err := os.Open("conf/configuration.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config2 models.DomoticaConfiguration
	json.Unmarshal([]byte(byteValue), &config2)
	json.Unmarshal([]byte(byteValue), &config)
	fmt.Println("Loaded config", config)
}

func (this *DomoticaController) Post() {

	//jsoninfo := this.GetString("jsoninfo")
	//fmt.Println("", jsoninfo)

	body := this.Ctx.Input.RequestBody
	var request models.DomoticaRequest
	json.Unmarshal(body, &request)

	fmt.Println("command ", request.Entity, request.Target, request.Cmd)

	entityConfiguration := config.Entities[request.Entity]

	command := entityConfiguration.Commands[request.Cmd]
	targetEnv := entityConfiguration.Env[request.Target]

	var newCommand = replace(command, targetEnv)
	newCommand = replace(newCommand, config.SystemParameters)
	fmt.Println("Executing", newCommand)

	_, err := http.Get(newCommand)

	var code = 200
	if err != nil {
		code = 404
		fmt.Println(err)
	}

	this.Ctx.ResponseWriter.WriteHeader(code)
	//this.Data["json"] = "{\"success\":\"ok\"}"
	//this.ServeJSON()
}

func replace(str string, values map[string]string) string {
	for k, v := range values {
		fmt.Println(k)
		fmt.Println(v)
		str = strings.Replace(str, "{"+k+"}", v, -1)
	}

	return str
}

func (this *DomoticaController) Get() {
	fmt.Println("received request")

	this.Data["json"] = "{\"tst\":1}"
	this.ServeJSON()
}
