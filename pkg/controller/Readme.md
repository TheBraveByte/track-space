### Package Controller API Endpoints
The following are the endpoints of the "package controller" web application:

`GET / - Home Page`

This endpoint renders the home page of the web application. It returns an HTML page with an authentication flag set to 0.

`GET /contact - Contact Page`

This endpoint renders the contact page of the web application. It returns an HTML page with a contact form.

`POST /contact - Post Contact Form`

This endpoint is used to submit the contact form. It accepts form data in the request body, validates it, and sends confirmation and help desk messages to the user and the support team, respectively.

`GET /signup - Sign Up Page`

This endpoint renders the sign-up page of the web application. It returns an HTML page with a sign-up form.


Note: All the above endpoints are implemented in the TrackSpace struct, which implements the repository pattern to access multiple packages at once, including app configuration and database collections.

Usage:
The above endpoints are intended to be used in conjunction with a web server, such as Gin. To use these endpoints, first create a new instance of the TrackSpace struct and initialize it with the app configuration and database client. Then, attach the desired endpoint functions to the appropriate routes in your web server's router.

Example:
```go
import (
    "github.com/gin-gonic/gin"
    "github.com/yusuf/track-space/pkg/controller"
    "github.com/yusuf/track-space/pkg/config"
    "go.mongodb.org/mongo-driver/mongo"
)

func main() {
    // Load app configuration
    appConfig := config.NewAppConfig()

    // Connect to the database
    dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.DBURI))
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Create a new instance of the TrackSpace struct
    ts := controller.NewTrackSpace(appConfig, dbClient)

    // Create a new router instance
    router := gin.Default()

    // Attach the endpoint functions to the appropriate routes
    router.GET("/", ts.HomePage())
    router.GET("/contact", ts.Contact())
    router.POST("/contact", ts.PostContact())
    router.GET("/signup", ts.SignUpPage())

    // Start the web server
    router.Run(":8080")
}
```

Note: The above example assumes that you have imported the necessary packages and that the required dependencies have been installed.
