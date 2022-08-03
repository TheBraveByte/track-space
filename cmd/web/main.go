package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"os"

	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/controller"
	"github.com/yusuf/track-space/pkg/db"
)

var app config.AppConfig

func main() {
	app.AppInProduction = false
	app.UseTempCache = false

	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file available")
	}

	mongodbUri := os.Getenv("MONGODB_URI")
	if mongodbUri == "" {
		log.Println("mongodb cluster uri not found : ")
	}

	portNumber := os.Getenv("PORT_NUMBER")
	if portNumber == "" {
		log.Println("No local server port number created!")
	}

	Client := db.DatabaseConnection(mongodbUri)
	// mail = data.MailData(Client, "mail")

	defer func() {
		if err = Client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
			return
		}
	}()

	repo := controller.NewTrackSpace(&app, Client)

	appRouter := gin.New()
	err = appRouter.SetTrustedProxies([]string{"127.0.0.1"})

	if err != nil {
		log.Println(err)
		log.Println("cannot access untrusted server proxy")

	}

	appRouter.SetFuncMap(template.FuncMap{})
	appRouter.Static("/static", "./static")
	appRouter.LoadHTMLGlob("templates/*.html")

	Routes(appRouter, *repo)

	err = appRouter.Run(portNumber)
	if err != nil {
		log.Fatal(err)
	}
}
