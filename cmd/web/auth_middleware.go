package main

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yusuf/track-space/pkg/auth"
	_ "log"
	"net/http"
	"os"
)

func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := sessions.Default(c)
		data.Get("session")
		authHeader := c.Writer.Header().Get("Authorization")
		authToken := authHeader[len(os.Getenv("TOKEN_SCHEME"))+1:]

		if authToken == "" {
			_ = c.AbortWithError(http.StatusNoContent, errors.New("no value for token"))
			return
		}

		authClaims, err := auth.ParseToken(authToken)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
			return
		}
		c.Set("email", authClaims.Email)
		c.Set("password", authClaims.Password)
		c.Set("uid", authClaims.IPAddress)
		c.Next()
	}

}
