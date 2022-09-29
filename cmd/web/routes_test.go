package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yusuf/track-space/pkg/controller"
	"testing"
)

func TestRoutes(t *testing.T) {
	router := gin.New()
	var repo controller.TrackSpace
	Routes(router, repo)
}
