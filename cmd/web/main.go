package main

import (
	"context"
	"encoding/gob"
	"html/template"
	"log"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/controller"
	"github.com/yusuf/track-space/pkg/driver"
	"github.com/yusuf/track-space/pkg/model"
	"github.com/yusuf/track-space/pkg/ws"
	"github.com/yusuf/track-space/pkg/wsconfig"
	"github.com/yusuf/track-space/pkg/wsmodel"
)

var (
	app      config.AppConfig
	validate *validator.Validate
)

func main() {
	gob.Register(model.User{})
	gob.Register(model.Auth{})
	gob.Register(model.Project{})
	gob.Register(model.Todo{})
	gob.Register(model.Email{})
	gob.Register(model.Data{})
	gob.Register(model.SessionData{})
	gob.Register(wsconfig.SocketConnection{})
	gob.Register(wsmodel.SocketPayLoad{})
	gob.Register(wsmodel.SocketResponse{})

	// Validate - to help check for a validated json database model
	validate = validator.New()

	mailChannel := make(chan model.Email)
	app.MailChan = mailChannel

	app.Validator = validate

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file available")
	}

	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatalln("mongodb cluster uri not found : ")
	}

	portNumber := os.Getenv("PORT_NUMBER")
	if portNumber == "" {
		log.Fatalln("No local server port number created!")
	}

	mailPass := os.Getenv("MAIL_PASSWORD")

	defer close(app.MailChan)

	log.Println("Application starting mail server listening to channel")
	// Listening to the localhost mail server
	go ListenToMailChannel(mailPass)

	// Listening to PayLoad from the websocket
	go ws.GetDataFromChannel()

	// connecting to the database
	Client := db.DatabaseConnection(mongodbURI)

	defer func() {
		if err := Client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
			return
		}
	}()

	repo := controller.NewTrackSpace(&app, Client)

	gin.SetMode(gin.ReleaseMode)
	appRouter := gin.New()
	proxyErr := appRouter.SetTrustedProxies([]string{"127.0.0.1"})

	if proxyErr != nil {
		log.Println(proxyErr)
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
