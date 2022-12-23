package main

import (
	"errors"
	"fmt"
	"net/http"
	_"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yusuf/track-space/pkg/auth"
)

// IsAuthorized Middleware for Authenticating the user
func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		token := tsData.Get("token").(string)
		if token == "" {
			_ = c.AbortWithError(http.StatusNoContent, errors.New("no value for token"))
			return
		}
		fmt.Println(token)

		authClaims, err := auth.ParseToken(token)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
			return
		}
		c.Set("token", token)
		c.Set("email", authClaims.Email)
		c.Set("_id", authClaims.ID)
		c.Set("uid", authClaims.IPAddress)
		c.Next()
	}
}
