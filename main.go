package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/giovannicammarata/gc-simple-home/configurator"
	"github.com/giovannicammarata/gc-simple-home/controllers"
	"github.com/giovannicammarata/gc-simple-home/models"
	"log"
	"os"
	"strings"
)

var config *models.SystemConfiguration
var configFile string

func main() {
	argsWithoutProg := os.Args[1:]

	configFile = "conf/configuration.json"
	if len(argsWithoutProg) > 0 {
		configFile = argsWithoutProg[0]
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

	loadConfiguration()

	go startFileWatcher()

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

func startFileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("error starting the watcher", err)
		os.Exit(2)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				loadConfiguration()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error in watcher:", err)
			}
		}
	}()

	err = watcher.Add(configFile)
	if err != nil {
		log.Fatal("error starting the watcher", err)
	}
	<-done
}

func loadConfiguration() {
	log.Println("loading the configuration from", configFile)
	// load configuration from file
	configRaw := models.SystemConfiguration{}
	err := configurator.LoadConfiguration(configFile, &configRaw)

	if err != nil && config != nil {
		log.Println("unable to load the configuration. Exiting")
		os.Exit(3)
	}

	m2 := make(map[string]models.EntityConfiguration)
	entities := configRaw.Domotica.Entities
	for k, v := range *entities {
		env := v.Env
		m := make(map[string]map[string]string)
		for k2, v2 := range env {
			keys := strings.Split(k2, "|")
			for _, key := range keys {
				m[key] = v2
			}
		}

		env = m
		m2[k] = models.EntityConfiguration{Env: env, Commands: v.Commands}
	}

	configRaw.Domotica.Entities = &m2
	config = &configRaw

	log.Println("loaded config")
	spew.Dump(config)
}
