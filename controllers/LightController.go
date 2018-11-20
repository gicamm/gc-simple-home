package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/giovannicammarata/simple_home/models"
)

type LightController struct {
	BaseController
}

func (this *LightController) Post() {

	//jsoninfo := this.GetString("jsoninfo")
	//fmt.Println("", jsoninfo)

	body := this.Ctx.Input.RequestBody
	var lightRequest models.LightRequest
	json.Unmarshal(body, &lightRequest)

	fmt.Println("Light command ", lightRequest.Cmd, lightRequest.Target)

	this.Ctx.ResponseWriter.WriteHeader(200)
	//this.Data["json"] = "{\"success\":\"ok\"}"
	//this.ServeJSON()
}

func (this *LightController) Get() {
	fmt.Println("received request")

	this.Data["json"] = "{\"tst\":1}"
	this.ServeJSON()
}
