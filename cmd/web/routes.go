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
	router.GET("/reset-password", h.ResetPassword())
	router.POST("/reset-password", h.UpdatePassword())

	authRouter := routes.Group("/auth")

	authRouter.Use(IsAuthorized())
	{
		// authRouter.Handle(http.MethodConnect, "/workspace", h.ProcessWorkSpace())
		authRouter.GET("/user/dashboard", h.GetDashBoard())

		authRouter.GET("/user/workspace", h.ProjectWorkspace())
		authRouter.POST("/user/workspace", h.PostWorkSpaceProject())

		authRouter.GET("/user/project-table", h.ShowProjectTable())
		authRouter.GET("/user/:src/:id/show-project", h.ShowUserProject())
		authRouter.POST("/user/project-table/:src/:id/change", h.ModifyUserProject())
		authRouter.GET("/user/:src/:id/delete", h.DeleteProject())

		authRouter.GET("/user/todo", h.GetTodo())
		authRouter.POST("/user/todo", h.PostTodoData())

		authRouter.GET("/user/todo-table", h.ShowTodoTable())
		authRouter.GET("/user/:src/:id/show-todo", h.ShowTodoSchedule())
		authRouter.POST("/user/todo-table/:src/:id/change", h.ModifyUserTodo())
		authRouter.GET("/user/todo/:src/:id/delete", h.DeleteTodo())

		// Routes for websocket handlers
		authRouter.GET("/user/chat", h.ChatRoom())
		authRouter.GET("/ts", h.ChatRoomEndpoint())
		authRouter.POST("/user/logout")

		//Admin routes
		authRouter.GET("/admin", h.AdminPage())
		authRouter.GET("/:src/dashboard/:id/delete", h.AdminDeleteUser())

	}
}
