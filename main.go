package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/davecgh/go-spew/spew"
	"github.com/giovannicammarata/simple_home/configurator"
	"github.com/giovannicammarata/simple_home/controllers"
	"github.com/giovannicammarata/simple_home/models"
	"os"
	"strings"
)

var config models.SystemConfiguration

func main() {
	argsWithoutProg := os.Args[1:]

	confFile := "conf/configuration.json"
	if len(argsWithoutProg) > 0 {
		confFile = argsWithoutProg[0]
	}

	// Authorization filter
	var authFilter = func(ctx *context.Context) {
		token := ctx.Input.Header("Authorization")
		if config.Network.Token == "" || strings.Compare(config.Network.Token, token) == 0 {
			return
		}

		ctx.Abort(404, "")
	}

	// setup the beego router
	beego.InsertFilter("/*", beego.BeforeRouter, authFilter)
	beego.Router("/", &controllers.DomoticaController{})

	// load configuration from file
	config = models.SystemConfiguration{}
	configurator.LoadConfiguration(confFile, &config)
	fmt.Println("Loaded config")
	spew.Dump(config)

	// Enable HTTP/HTTPs
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

	beego.Run() // start the router
}
