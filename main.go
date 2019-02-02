package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/davecgh/go-spew/spew"
	"github.com/giovannicammarata/simple_home/configurator"
	"github.com/giovannicammarata/simple_home/controllers"
	"github.com/giovannicammarata/simple_home/models"
	"os"
)

func main() {
	//beego.Router("/light", &controllers.DomoticaController{}, "post:Post", "get:Get")

	argsWithoutProg := os.Args[1:]

	confFile := "conf/configuration.json"
	if len(argsWithoutProg) > 0 {
		confFile = argsWithoutProg[0]
	}

	beego.Router("/", &controllers.DomoticaController{})

	config := models.SystemConfiguration{}
	configurator.LoadConfiguration(confFile, &config)
	fmt.Println("Loaded config")
	spew.Dump(config)

	if config.Network.HTTPPort > 0 {
		beego.BConfig.Listen.HTTPPort = config.Network.HTTPPort
	}
	if config.Network.HTTPSPort > 0 {
		beego.BConfig.Listen.EnableHTTPS = true
		beego.BConfig.Listen.HTTPSPort = config.Network.HTTPSPort
		beego.BConfig.Listen.HTTPSCertFile = config.Network.HTTPSCertFile
		beego.BConfig.Listen.HTTPSKeyFile = config.Network.HTTPSKeyFile
	}
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

	controllers.Config = &config.Domotica

	beego.Run()
}
