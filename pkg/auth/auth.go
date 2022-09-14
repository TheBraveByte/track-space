package auth

import (
	// "fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TrackClaims type struct which is used to create / generate jwt token
type TrackClaims struct {
	jwt.StandardClaims
	Email     string
	ID        string
	IPAddress string
}

// GenerateJWTToken : This functions helps to create a JWT token using the
// SignedStringMethod of the ES256 algorithm using a TOKEN_KEY and the claims
// to generate a token
func GenerateJWTToken(email, id, ipaddress string) (string, string, error) {
	trackToken := TrackClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Local().Add(time.Duration(12045) * time.Hour).Unix(),
			Issuer:    "trackSpace",
		},
		email,
		id,
		ipaddress,
	}
	refreshToken := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Duration(485994) * time.Hour).Unix(),
		Issuer:    "trackSpace",
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, trackToken).SignedString([]byte(os.Getenv("TOKEN")))
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	newToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshToken).SignedString([]byte(os.Getenv("TOKEN")))
	if err != nil {
		log.Println(err, "dbj ")
		return "", "", err
	}
	return token, newToken, nil
}

// ParseToken : this function helps to validate the generated JSON WEB TOKEN(JWT)
// and also check for errors
func ParseToken(tokenValue string) (*TrackClaims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &TrackClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN")), nil
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	tokenClaim, ok := token.Claims.(*TrackClaims)
	if !ok {
		log.Println(http.StatusUnauthorized, "Invalid token claim")
	}
	err = tokenClaim.Valid()
	if err != nil {
		log.Println(http.StatusUnauthorized, "Generate token invalid")
	}
	return tokenClaim, nil
}
