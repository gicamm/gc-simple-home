package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/giovannicammarata/simple_home/controllers"
)

func main() {
	//beego.Router("/light", &controllers.DomoticaController{}, "post:Post", "get:Get")
	beego.Router("/", &controllers.DomoticaController{})

	beego.BConfig.Listen.EnableHTTPS = true
	beego.BConfig.Listen.HTTPPort = 60001
	beego.BConfig.Listen.HTTPSPort = 60002
	beego.BConfig.Listen.HTTPSCertFile = "conf/host.crt"
	beego.BConfig.Listen.HTTPSKeyFile = "conf/host.key"
	beego.BConfig.CopyRequestBody = true

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()
}
