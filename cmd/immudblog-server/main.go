package main

import (
	"fmt"
	"immudblog/config"
	"immudblog/immudb"
	"immudblog/restapi"
	"os"

	log "github.com/sirupsen/logrus"
	ginFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"immudblog/docs"
)

// @title           Immudb Logstore
// @version         1.0
// @description     A sample application for store loglines to immudb

// @contact.name   Akos Balogh
// @contact.email  akosbalogh005@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic	BasicAuth

// @host      localhost:8080
// @BasePath  /api/v1
func main() {

	config.Init()
	log.SetOutput(os.Stdout)

	// setup log
	if config.LogFlags.Type == config.LOG_TYPE_JSON {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetLevel(log.InfoLevel)
	config.LogFlags.Debug = true
	if config.LogFlags.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Infof("Starting immudblog application...")
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Immudb Logstore API"
	docs.SwaggerInfo.Description = "A sample application for store loglines to immudb"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", config.ServerFlags.Host, config.ServerFlags.Port)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := restapi.SetupRouter()

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(ginFiles.Handler))

	immudb.Init()

	listenon := fmt.Sprintf(":%d", config.ServerFlags.Port)
	log.Infof("Starting immudblog REST API (%v)...", listenon)
	r.Run(listenon)

}
