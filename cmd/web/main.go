package main

import (
	"context"
	"encoding/gob"
	"html/template"
	"log"
	"os"
	


	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/controller"
	"github.com/yusuf/track-space/pkg/db"
	"github.com/yusuf/track-space/pkg/model"
	"github.com/yusuf/track-space/pkg/ws"
)

var app config.AppConfig

func main() {
	gob.Register(model.User{})
	gob.Register(model.Auth{})
	gob.Register(model.Project{})
	gob.Register(model.DailyTask{})
	gob.Register(model.Email{})
	app.AppInProduction = false
	app.UseTempCache = false

	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file available")
	}

	mongodbUri := os.Getenv("MONGODB_URI")
	if mongodbUri == "" {
		log.Println("mongodb cluster uri not found : ")
		return
	}

	portNumber := os.Getenv("PORT_NUMBER")
	if portNumber == "" {
		log.Println("No local server port number created!")
		return
	}
	go ws.GetDataFromChannel()

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
	appRouter.Static("/main-ckeditor5", "./main-ckeditor5")
	// appRouter.StaticFS("./static/", http.Dir("static"))
	// appRouter.StaticFile("/favicon.ico", "./resources/favicon.ico")
	appRouter.LoadHTMLGlob("templates/*.html")

	Routes(appRouter, *repo)

	err = appRouter.Run(portNumber)
	if err != nil {
		log.Fatal(err)
	}
}
