package controller

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/model"
	"github.com/yusuf/track-space/pkg/ws"
	"github.com/yusuf/track-space/pkg/wsconfig"
	"github.com/yusuf/track-space/pkg/wsmodel"
	"log"
	"text/template"
)

var appConfig config.AppConfig

// Creating a mock testing for controller package
func TrackSpaceSetUp() *gin.Engine {
	gob.Register(model.User{})
	gob.Register(model.Auth{})
	gob.Register(model.Project{})
	gob.Register(model.Todo{})
	gob.Register(model.Email{})
	gob.Register(model.Data{})
	gob.Register(wsconfig.SocketConnection{})
	gob.Register(wsmodel.SocketPayLoad{})
	gob.Register(wsmodel.SocketResponse{})

	mailChannel := make(chan model.Email)
	appConfig.MailChan = mailChannel
	appConfig.AppInProduction = false
	appConfig.UseTempCache = false

	defer close(appConfig.MailChan)

	log.Println("Application starting mail server listening to channel")
	// Listening to the localhost mail server

	//Listening to PayLoad from the websocket
	go ws.GetDataFromChannel()

	appRouter := gin.New()
	err := appRouter.SetTrustedProxies([]string{"127.0.0.1"})

	if err != nil {
		log.Println(err)
		log.Println("cannot access untrusted server proxy")

	}

	appRouter.SetFuncMap(template.FuncMap{})
	appRouter.Static("/static", "./static")
	appRouter.LoadHTMLGlob("./../../templates/*.html")

	appRouter.Use(gin.Logger(), gin.Recovery())
	storeData := cookie.NewStore([]byte("trackSpace"))
	appRouter.Use(sessions.Sessions("session", storeData))
	return appRouter
}
