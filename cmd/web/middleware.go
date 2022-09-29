package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yusuf/track-space/pkg/auth"
)

// IsAuthorized Middleware for Authenticating the user
func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := c.Request.Cookie("bearerToken")
		if err != nil {
			if err == http.ErrNoCookie {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		if authToken.Value == "" {
			_ = c.AbortWithError(http.StatusNoContent, errors.New("no value for token"))
			return
		}

		authClaims, err := auth.ParseToken(authToken.Value)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
			return
		}
		c.Set("token", authToken.Value)
		c.Set("email", authClaims.Email)
		c.Set("_id", authClaims.ID)
		c.Set("uid", authClaims.IPAddress)
		c.Next()
	}
}
