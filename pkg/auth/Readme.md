# Auth Package Documentation
The auth package provides functions to generate and validate JSON Web Tokens (JWTs). The package uses the Go standard library's net/http and log packages as well as the third-party package github.com/golang-jwt/jwt/v4.

### Installation
To install this package, use the following command:

`go get github.com/<username>/<reponame>/auth`
  
### Usage
To use this package, import it into your Go code as follows:

```go
import "github.com/<username>/<reponame>/auth"
```

### Features
The auth package provides the following features:

### GenerateJWTToken
The GenerateJWTToken function creates a JWT token using the SignedStringMethod of the ES256 algorithm. The function takes three arguments: email, ID, and IP address. The function returns the JWT token and a refresh token. The refresh token has the same expiration time as the JWT token.

```go
func GenerateJWTToken(email, id, ipaddress string) (string, string, error)
```

### ParseToken

The ParseToken function validates a JWT token and checks for errors. The function takes a JWT token as an argument and returns a TrackClaims struct and an error. The TrackClaims struct contains the claims used to generate the token.

```go

func ParseToken(tokenValue string) (*TrackClaims, error)
```

Example
The following example shows how to generate a JWT token and parse it using the auth package.

```go

package main

import (
	"fmt"
	"log"

	"github.com/<username>/<reponame>/auth"
)

func main() {
	// Generate a JWT token
	token, refreshToken, err := auth.GenerateJWTToken("johndoe@example.com", "123456", "127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JWT token
	claims, err := auth.ParseToken(token)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("JWT Token:", token)
	fmt.Println("Refresh Token:", refreshToken)
	fmt.Println("Email:", claims.Email)
	fmt.Println("ID:", claims.ID)
	fmt.Println("IP Address:", claims.IPAddress)
}
```
