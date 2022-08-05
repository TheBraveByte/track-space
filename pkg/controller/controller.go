package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/yusuf/track-space/pkg/data"
	"github.com/yusuf/track-space/pkg/data/tsRepoStore"
	"github.com/yusuf/track-space/pkg/key"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yusuf/track-space/pkg/auth"
	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/model"
)

var Validate = validator.New()

// TrackSpace Implement the repository pattern to access multiple package all at once
// this will give me access to the app configuration package
// and the database collections as well
type TrackSpace struct {
	AppConfig *config.AppConfig
	tsDB      data.TrackSpaceDBRepo
}

func NewTrackSpace(appConfig *config.AppConfig, tsm *mongo.Client) *TrackSpace {
	return &TrackSpace{
		AppConfig: appConfig,
		tsDB:      tsRepoStore.NewTsMongoDBRepo(appConfig, tsm),
	}
}

func (ts *TrackSpace) HomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home-page.html", gin.H{})
	}
}

func (ts *TrackSpace) SignUpPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup-page.html", gin.H{})
	}
}

func (ts *TrackSpace) PostSignUpPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User

		if err := Validate.Struct(&user); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		email := c.PostForm("email")
		password := key.HashPassword(c.PostForm("password"))

		tsData := sessions.Default(c)
		tsData.Set("email", email)
		tsData.Set("password", password)

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}

		_, err := ts.tsDB.InsertInfo(email, password)

		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		if err == nil {
			c.Redirect(http.StatusSeeOther, "/")
		}

		c.Redirect(http.StatusSeeOther, "/user-info")

	}
}

func (ts *TrackSpace) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-info.html", gin.H{})
	}
}

func (ts *TrackSpace) PostUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		tsData := sessions.Default(c)

		validateErr := Validate.Struct(&user)
		if validateErr != nil {
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: validateErr})
			return
		}

		err := c.Request.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}

		info := make(map[string]interface{})
		email := tsData.Get("email")

		info["userID"] = primitive.NewObjectID()
		info["firstname"] = c.Request.Form.Get("first-name")
		info["lastname"] = c.Request.Form.Get("last-name")
		info["YrsOfExp"] = c.Request.Form.Get("yrs-of-exp")
		info["country"] = c.Request.Form.Get("nation")
		info["phone"] = c.Request.Form.Get("phone")

		info["created_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		info["updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		info["stack"] = append(user.Stack, c.Request.Form.Get("stack-name"))

		info["address"] = c.Request.Form.Get("address")

		var t1, t2 = "", ""

		err = ts.tsDB.UpdateUserInfo(info, email, t1, t2)
		if err != nil {
			log.Println("Cannot update user info")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		c.Redirect(http.StatusSeeOther, "/login")

	}
}

func (ts *TrackSpace) GetLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login-page.html", gin.H{})
	}
}

func (ts *TrackSpace) PostLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)

		if err := c.Request.ParseForm(); err != nil {
			log.Println("error while parsing form")
			return
		}

		email := tsData.Get("email")
		password := tsData.Get("password")

		IPAddress := c.Request.RemoteAddr

		//Posted form value
		postEmail := c.PostForm("email")
		postPassword := c.PostForm("password")

		//Previous store details in the database
		ok, hashedPassword := ts.tsDB.VerifyLogin(postEmail)

		if ok {
			ok, msg := key.VerifyPassword(postPassword, hashedPassword)
			log.Println(msg)
			if postPassword == password && postEmail == email && ok == true {
				token, newToken, err := auth.GenerateJWTToken(email, password, IPAddress)

				if err != nil {
					log.Println("cannot generate json web token")
					_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
					return
				}
				info := make(map[string]interface{})

				authData := make(map[string][]string)
				authData["auth"] = []string{token, newToken}
				t1 := authData["auth"][0]
				t2 := authData["auth"][1]
				fmt.Println(t1, t2)
				err = ts.tsDB.UpdateUserInfo(info, email, t1, t2)
				if err != nil {
					log.Println("cannot update user info")
					return
				}
				c.Writer.Header().Set("Authorization", fmt.Sprintf("BearerToken %s", t1))
				tsData.AddFlash("Successfully login")

			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
				tsData.AddFlash("Incorrect Email or Password")
				c.Abort()
				return
			}
		}

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}

		//c.JSON(http.StatusOK, "Successfully use token")
		c.HTML(http.StatusOK, "/", gin.H{})

	}
}


//GetDashBoard :: a lot of logic will be done here ..... alot
func (ts *TrackSpace) GetDashBoard() gin.HandlerFunc {
	return func(c *gin.Context) {

		tsData := sessions.Default(c)
		email := tsData.Get("email")

		user, err := ts.tsDB.SendUserDetails(email)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}

		//this controller with still be uppdated as i progress

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}
		c.HTML(http.StatusOK, "dash.html", gin.H{
			"FirstName": user["first_name"],
			"LastName":  user["last_name"],
		})
	}
}

var StartTime time.Time
func (ts *TrackSpace) WorkSpace() gin.HandlerFunc{
	return func(c *gin.Context){
		
		c.HTML(http.StatusOK,"work.html",gin.H{
			"StartTime": StartTime.Local().UTC(),
		})
	}
}


func(ts *TrackSpace) PostWorkSpace() gin.HandlerFunc{
	return func(c *gin.Context){
		var project model.Project
		// var EndTime time.Time
		tsData := sessions.Default(c)
		userEmail  := tsData.Get("email")

		if validateErr := Validate.Struct(&project); validateErr != nil{
			log.Println("cannot validate project struct")
			_ =c.AbortWithError(http.StatusBadRequest,gin.Error{Err: validateErr})
			return
		}

		//Getting the project data
		if err:= c.Request.ParseForm() ; err != nil {
			log.Println("cannot parse the workspace form")
			return
		}

		projectData := make(map[string]interface{})
		projectData["project_name"] = c.PostForm("project-name")
		projectData["tools_use_as"] = c.PostForm("project-tool-use")
		projectData["project_content"] = c.PostForm("editor")
		projectData["start_time"] =  StartTime.Local().UTC()
		err := ts.tsDB.StoreWorkSpaceData(userEmail, projectData)
		if err != nil{
			log.Println("Error while storing using user project data")
			return
		}
		tsData.AddFlash("successfully submitted project")

		c.HTML(http.StatusOK,"work.html",gin.H{
			"Save":tsData.Flashes("successfully submitted project"),
		})
	}
}
