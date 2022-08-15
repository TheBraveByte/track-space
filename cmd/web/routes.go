package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/yusuf/track-space/pkg/controller"
)

func Routes(routes *gin.Engine, h controller.TrackSpace) {
	router := routes.Use(gin.Logger(), gin.Recovery())
	router.GET("/", h.HomePage())

	storeData := cookie.NewStore([]byte("trackSpace"))
	router.Use(sessions.Sessions("session", storeData))
	router.GET("/sign-up", h.SignUpPage())
	router.POST("/sign-up", h.PostSignUpPage())

	router.GET("/user-info", h.GetUserInfo())
	router.POST("/user-info", h.PostUserInfo())

	router.GET("/login", h.GetLoginPage())
	router.POST("/login", h.PostLoginPage())

	authRouter := routes.Group("/auth")
	authRouter.Use(IsAuthorized())
	{
		// authRouter.Handle(http.MethodConnect, "/workspace", h.ProcessWorkSpace())
		authRouter.GET("/user/dashboard", h.GetDashBoard())
		authRouter.GET("/user/workspace", h.WorkSpace())
		authRouter.POST("/user/workspace/save", h.PostWorkSpace())
		authRouter.GET("user/daily-task", h.DailyTaskTodo())
		authRouter.POST("user/daily-task/save", h.PostDailyTaskTodo())
		authRouter.GET("user/orkspace/:src/:id/show-project", h.ShowUserProject())
		authRouter.POST("/user/logout")
	}
}
