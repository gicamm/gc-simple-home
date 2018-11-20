package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/giovannicammarata/simple_home/controllers"
)

func main() {
	//beego.Router("/light", &controllers.LightController{}, "post:Post", "get:Get")
	beego.Router("/light", &controllers.LightController{})

	beego.BConfig.Listen.EnableHTTPS = true
	beego.BConfig.Listen.HTTPSPort = 10443
	beego.BConfig.Listen.HTTPSCertFile = "/home/sparrow/workspace/go/src/github.com/giovannicammarata/simple_home/resources/conf/host.crt"
	beego.BConfig.Listen.HTTPSKeyFile = "/home/sparrow/workspace/go/src/github.com/giovannicammarata/simple_home/resources/conf/host.key"
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
