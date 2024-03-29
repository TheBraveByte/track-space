package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/yusuf/track-space/pkg/data"
	"github.com/yusuf/track-space/pkg/data/tsRepoStore"
	"github.com/yusuf/track-space/pkg/key"
	"github.com/yusuf/track-space/pkg/temp"
	"github.com/yusuf/track-space/pkg/ws"
	"github.com/yusuf/track-space/pkg/wsconfig"
	"github.com/yusuf/track-space/pkg/wsmodel"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yusuf/track-space/pkg/auth"
	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/model"
)

/*
TrackSpace Implement the repository pattern to access multiple package
all at once this will give me access to the app configuration package
and the database collections as well
*/
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

// NewTestTrackSpace A dump copy of the tracks-pace struct for unit testing
func NewTestTrackSpace(appConfig *config.AppConfig) *TrackSpace {
	return &TrackSpace{
		AppConfig: appConfig,
		tsDB:      tsRepoStore.NewTsMongoDBRepo(appConfig, nil),
	}
}

func (ts *TrackSpace) HomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home-page.html", gin.H{
			"authenticate": 0,
		})
	}
}

func (ts *TrackSpace) Contact() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.html", gin.H{})
	}
}

func (ts *TrackSpace) PostContact() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			log.Panic("form not parsed")
			return
		}
		name := c.Request.Form.Get("full-name")
		email := c.Request.Form.Get("email")
		msg := c.Request.Form.Get("message")

		TeamMessage := fmt.Sprintf(`
			<strong>Help Desk Message</strong><br>
			<p>Email from %s</p>
			Hi, %s:<br>
			<p>%s</p>
			`, email ,"Trackspace Team", msg)
		TeamMailMsg := model.Email{
			Subject:  "Help Desk Message",
			Content:  TeamMessage,
			Sender:   "official.trackspace@gmail.com",
			Receiver: "official.trackspace@gmail.com",
			Template: "email.html",
		}
		ts.AppConfig.MailChan <- TeamMailMsg

		message := fmt.Sprintf(`
			<strong>Confirmation Message</strong><br>
			Hi, %s <br>
			<p>This is to confirm that your message have been received 
            by track-space team. You will hear from us in few days time.
			Feel free to explore our core service and other features
			</p>
			`, name)
		mailMsg := model.Email{
			Subject:  "Confirmation Message",
			Content:  message,
			Sender:   "official.trackspace@gmail.com",
			Receiver: email,
			Template: "email.html",
		}

		ts.AppConfig.MailChan <- mailMsg

		c.HTML(http.StatusOK, "contact.html", gin.H{})
	}
}

// SignUpPage - Handler to get the sign-up page for user
func (ts *TrackSpace) SignUpPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup-page.html", gin.H{})
	}
}

/*
PostSignUpPage - this validates the user input and store the value in
session as cookies for future usage, insert user input in the database
check for existing user and hashed the user password as well
*/
func (ts *TrackSpace) PostSignUpPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User

		if err := c.Request.ParseForm(); err != nil {
			log.Panic("form not parsed")
			return
		}
		user.Email = c.Request.Form.Get("email")
		user.Password = key.HashPassword(c.Request.Form.Get("password"))

		// Server side validation of the user input from a form
		if err := ts.AppConfig.Validator.Struct(user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}

		count, userID, err := ts.tsDB.InsertUserInfo(user.Email, user.Password)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}
		tsData := sessions.Default(c)
		userData := model.SessionData{
			UserID:   userID,
			Email:    user.Email,
			Password: user.Password,
		}
		tsData.Set("session_data", userData)

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}

		if count == 1 {
			c.HTML(http.StatusSeeOther, "login-page.html", gin.H{
				"msg": "Email already registered on track-space. Log-in into your account",
			})
		} else {
			c.Redirect(http.StatusSeeOther, "/user-info")
		}
	}
}

// GetUserInfo - handler to get the user-details/ info page
func (ts *TrackSpace) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-info.html", gin.H{})
	}
}

/*
PostUserInfo - this validates the user model and get the user input details
and store the details in the database and this helps to redirect
the user to a login page
*/
func (ts *TrackSpace) PostUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User

		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)

		if err := c.Request.ParseForm(); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		user.FirstName = c.Request.Form.Get("first-name")
		user.LastName = c.Request.Form.Get("last-name")
		user.Address = c.Request.Form.Get("address")
		user.YrsOfExp = c.Request.Form.Get("yrs-of-exp")
		user.Profession = c.Request.Form.Get("profession")
		user.Country = c.Request.Form.Get("nation")
		user.PhoneNumber = c.Request.Form.Get("phone")
		user.Stack = append(user.Stack, c.Request.Form.Get("stack-name"))
		user.IPAddress = c.Request.RemoteAddr
		user.CreatedAt = time.Now().Format("2006-01-02")
		user.UpdatedAt = time.Now().Format("2006-01-02")
		t1, t2 := "", ""
		tsData.Set("first-name", user.FirstName)

		// Server side validation of the user input from a form
		if err := ts.AppConfig.Validator.Struct(user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				return
			}
		}
		err := ts.tsDB.UpdateUserInfo(user, userData.UserID, t1, t2)
		if err != nil {
			log.Println("Cannot update user info")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}
		message := fmt.Sprintf(`
			<strong>Confirmation for Account Created</strong><br>
			Hi, %s:<br>
			<p>This is to confirm that your have sign up for track-space.
			We hope you have a wonderful experience using our platform as
			your workspace for your project
			</p>
			`, user.FirstName)
		mailMsg := model.Email{
			Subject:  "Confirmation for Account Created",
			Content:  message,
			Sender:   "official.trackspace@gmail.com",
			Receiver: fmt.Sprint(tsData.Get("email")),
			Template: "email.html",
		}

		ts.AppConfig.MailChan <- mailMsg

		TeamMessage := fmt.Sprintf(`
			<strong>New User Account Notification</strong><br>
			Hi, %s:<br>
			<p>This is notify you guys that new user with an 
            <strong> ID:</strong> %s and <strong>IPAddress :</strong> of %s
            sign up for track-space.
			</p>
			`, "track-space Team", userData.UserID, user.IPAddress)
		TeamMailMsg := model.Email{
			Subject:  "Confirmation for Account Created",
			Content:  TeamMessage,
			Sender:   "official.trackspace@gmail.com",
			Receiver: "official.trackspace@gmail.com",
			Template: "email.html",
		}

		ts.AppConfig.MailChan <- TeamMailMsg
		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}
		c.Redirect(http.StatusSeeOther, "/login")
	}
}

// GetLoginPage - this handlers get the user login-page
func (ts *TrackSpace) GetLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login-page.html", gin.H{})
	}
}

/*
PostLoginPage : this handler help to verify the user password, authenticate other
user login details with respect to the database,generate an authorization token
for the user, as well as authorize the user and set the Response Header
with the Bearer Token
*/
func (ts *TrackSpace) PostLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)
		var templateData temp.TemplateData
		templateData.IsAuthenticated = 1

		if err := c.Request.ParseForm(); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		IPAddress := c.Request.RemoteAddr

		// Posted form value
		var user model.User
		user.Email = c.Request.Form.Get("email")
		user.Password = c.Request.Form.Get("password")

		// Server side validation of the user input from a form
		if err := ts.AppConfig.Validator.Struct(&user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}

		switch {
		case userData.Email == user.Email:
			// check to verify for the stored hashed password in database
			ok, _ := ts.tsDB.VerifyLogin(userData.UserID, userData.Password, user.Password)
			if ok {
				// check to match hashed password and the user password input
				token, newToken, err := auth.GenerateJWTToken(user.Email, userData.UserID, IPAddress)
				if err != nil {
					log.Println("cannot generate json web token")
					_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
					return
				}

				authData := make(map[string][]string)
				authData["auth"] = []string{token, newToken}
				tokenGen := authData["auth"][0]
				newTokenGen := authData["auth"][1]

				// saving the refresh token
				tsData.Set("token", tokenGen)
				tsData.Set("refreshToken", newTokenGen)
				if err := tsData.Save(); err != nil {
					log.Println("error from the session storage")
					_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
					return
				}

				err = ts.tsDB.UpdateUserField(userData.UserID, tokenGen, newTokenGen)
				if err != nil {
					_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
					return
				}

				c.HTML(http.StatusOK, "home-page.html", gin.H{
					"success":      "You have successfully logged-in on track-space. Go to dashboard",
					"authenticate": templateData.IsAuthenticated,
				})

			} else {
				c.HTML(http.StatusOK, "home-page.html", gin.H{
					"error": "invalid password, input correct password",
				})
			}
		case user.Email == "official.trackspace@gmail.com" && user.Password == "@_trackspace_":
			// Setting up the login authentication for admin
			adminInfo, err := ts.tsDB.GetAdminInfo()
			var (
				adminID       string
				adminPassword string
				adminEmail    string
			)
			for _, r := range adminInfo {
				for x, y := range r {
					if x == "_id" {
						adminID = fmt.Sprint(y)
					}
					if x == "email" {
						adminEmail = fmt.Sprint(y)
					}
					if x == "password" {
						adminPassword = fmt.Sprint(y)
					}
				}
			}
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}
			// check to verify for the stored hashed password in database
			ok, msg := key.VerifyPassword(user.Password, adminPassword)
			if !ok {
				log.Fatalf("Admin -- %s", msg)
				return
			}
			adminIPAddress := c.Request.RemoteAddr

			token, newToken, err := auth.GenerateJWTToken(adminEmail, adminID, adminIPAddress)
			if err != nil {
				log.Println("cannot generate json web token")
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}

			authData := make(map[string][]string)
			authData["auth"] = []string{token, newToken}
			tokenGen := authData["auth"][0]
			newTokenGen := authData["auth"][1]

			tsData.Set("refreshToken", newTokenGen)
			err = ts.tsDB.UpdateUserField(userData.UserID, tokenGen, newTokenGen)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}

			// c.Header("Authorization", "Bearer "+tokenGen)
			// c.SetCookie("bearerToken", tokenGen, 60*60*24*1200, "/", "localhost", false, true)
			// log.Println("Successfully login")

			c.HTML(http.StatusOK, "home-page.html", gin.H{
				"success":   "logged in successfully! Go to Admin",
				"authAdmin": templateData.IsAuthenticated,
			})
		default:
			c.HTML(http.StatusNotFound, "home-page.html", gin.H{
				"error": "incorrect password and email, Sign up your account on track space here!",
			})
		}
	}
}

/*
ResetPassword : this will help user to changes their previous password to a
new changes when forgotten
*/
func (ts *TrackSpace) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "reset.html", gin.H{})
	}
}

func (ts *TrackSpace) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		tsData := sessions.Default(c)

		password := fmt.Sprint(tsData.Get("password"))
		if err := ts.AppConfig.Validator.Struct(&user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			}
		}

		if err := c.Request.ParseForm(); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}
		user.Email = c.Request.Form.Get("user-email")
		user.Password = key.HashPassword(c.Request.Form.Get("new-password"))

		if user.Password != password {
			err := ts.tsDB.ResetUserPassword(user.Email, user.Password)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}
			TeamMessage := fmt.Sprintf(`
			<strong>Reset User Password</strong><br>
			Hi, %s:<br>
            <p>This is notify the team that user with an 
			<strong> ID : </strong> %s reset account password.
			</p>
			`, "Track-space Team", tsData.Get("userID"))
			TeamMailMsg := model.Email{
				Subject:  "Password Reset",
				Content:  TeamMessage,
				Sender:   "official.trackspace@gmail.com",
				Receiver: "official.trackspace@gmail.com",
				Template: "email.html",
			}

			ts.AppConfig.MailChan <- TeamMailMsg

			c.HTML(http.StatusSeeOther, "login-page.html", gin.H{
				"resetMsg": "password successfully reset. Log-in into your account",
			})
		} else {
			c.HTML(http.StatusTemporaryRedirect, "login-page.html", gin.H{
				"resetMsg": "Cannot reset password. New Password same as previous",
			})
		}
	}
}

// Todo Each of the  functions are use to organize the user data in the database
// Todo
func (ts *TrackSpace) Todo(count map[string]int, countTodo int) {
	count["todoNo"] = countTodo
}

func (ts *TrackSpace) Text(count map[string]int, countText int) {
	count["text"] = countText
}

func (ts *TrackSpace) Code(count map[string]int, countCode int) {
	count["code"] = countCode
}

func (ts *TrackSpace) Article(count map[string]int, countArticle int) {
	count["article"] = countArticle
}

// GetDashBoard - this show the user dashboard with respect to all the database
// details and queries; full brief or user activities
func (ts *TrackSpace) GetDashBoard() gin.HandlerFunc {
	return func(c *gin.Context) {
		t, ok := c.Get("token")
		if ok {
			tsData := sessions.Default(c)
			userData := tsData.Get("session_data").(model.SessionData)
			user, err := ts.tsDB.SendUserDetails(userData.UserID)
			if err != nil {
				_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
				return
			}
			// Count different projects type
			currentDate := time.Now().Format("2006-01-02")
			var storedDate string
			count := make(map[string]int)
			for k, value := range user {
				if k == "project_details" {
					switch v := value.(type) {
					case primitive.A:
						countCode, countText, countArticle := 0, 0, 0
						// _ is the index and y is the array of structs
						for _, y := range v {
							switch tools := y.(type) {
							case primitive.M:
								for i, j := range tools {
									// fmt.Println(i, j)
									if i == "created_at" {
										storedDate = fmt.Sprint(j)
									}
									if i == "tools_use_as" && j == "code" {
										countCode += 1
									} else if i == "tools_use_as" && j == "text" {
										countText += 1
									} else if i == "tools_use_as" && j == "article" {
										countArticle += 1
									}

								}
								ts.Code(count, countCode)
								ts.Text(count, countText)
								ts.Article(count, countArticle)
							}
						}
					}
				}
				if k == "todo" {
					switch todoList := value.(type) {
					case primitive.A:
						todoNo := len(todoList)
						ts.Todo(count, todoNo)
					}
				}
			}

			tsStat := model.Data{
				Date:    currentDate,
				Code:    count["code"],
				Article: count["article"],
				Text:    count["text"],
				Todo:    count["todoNo"],
				Total:   count["article"] + count["text"] + count["code"],
			}

			if currentDate == storedDate {
				err = ts.tsDB.UpdateUserStat(tsStat, userData.UserID)
				if err != nil {
					_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
					return
				}
			}
			r, err := ts.tsDB.GetUserStatByID(userData.UserID)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}
			statFile, err := json.MarshalIndent(r["data"], "", " ")
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}

			_ = ioutil.WriteFile("./static/json/data.json", statFile, 0o644)

			// this controller with still be updated as I progress
			if err := tsData.Save(); err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
				return
			}
			c.HTML(http.StatusOK, "dash.html", gin.H{
				"FirstName": user["first_name"],
				"LastName":  user["last_name"],
				"token":     t,
			})
		}
	}
}

/*
ProjectWorkspace :  this show the user workspace worksheet to execute

	projects and also to make use of other tools
*/
func (ts *TrackSpace) ProjectWorkspace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "work.html", gin.H{})
	}
}

/*
PostWorkSpaceProject : post the input content and store the content in the database of the
specific user
*/
func (ts *TrackSpace) PostWorkSpaceProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)

		if err := tsData.Save(); err != nil {
			_ = c.AbortWithError(http.StatusNoContent, gin.Error{Err: err})
			return
		}
		// Getting the project data
		if err := c.Request.ParseForm(); err != nil {
			log.Println("cannot parse the workspace form")
			return
		}
		project.ID = primitive.NewObjectID().Hex()
		project.ProjectName = strings.ToLower(c.PostForm("project-name"))
		project.ToolsUseAs = strings.ToLower(c.PostForm("project-tool-use"))
		project.ProjectContent = c.Request.Form.Get("myText")
		project.Status = "unmodified"
		project.CreatedAt = time.Now().Format("2006-01-02")
		project.UpdatedAt = time.Now().Format("2006-01-02")

		// Server side validation of the user input from a form
		if err := ts.AppConfig.Validator.Struct(&project); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}

		err := ts.tsDB.StoreProjectData(userData.UserID, project)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}

		// c.Redirect(http.StatusSeeOther, "/auth/user/workspace")
		c.HTML(http.StatusOK, "work.html", gin.H{
			"save": fmt.Sprintf("%v added to your projects", project.ProjectName),
		})
	}
}

/*
ShowProjectTable - this give the full projects details of a particular
user and help the user to modify each projects
*/
func (ts *TrackSpace) ShowProjectTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		project := make(map[string]interface{})
		var allProjects []map[string]interface{}

		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)

		user, err := ts.tsDB.SendUserDetails(userData.UserID)
		if err != nil {
			log.Println("cannot get user project data from the database")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}
		switch p := user["project_details"].(type) {
		case primitive.A:
			for _, x := range p {
				switch k := x.(type) {
				case primitive.M:
					for i, j := range k {
						project[i] = j
					}
					allProjects = append(allProjects, k)
				}
			}
		}
		c.HTML(http.StatusOK, "project-table.html", gin.H{
			"Project":   allProjects,
			"FirstName": user["first_name"],
			"LastName":  user["last_name"],
		})
	}
}

/*
ShowUserProject : this  handler direct the user to a page to make changes and modify their
existing projects store in the database
*/
func (ts *TrackSpace) ShowUserProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		projectMap := make(map[string]string)
		sourceLink := c.Param("src")
		project.ID = c.Param("id")
		ok := primitive.IsValidObjectID(c.Param("id"))

		if sourceLink != "project-table" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		projectData, err := ts.tsDB.GetProjectData(project.ID)
		if err != nil {
			log.Println("cannot get user project data from the database")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}

		switch p := projectData["project_details"].(type) {
		case primitive.A:
			for _, x := range p {
				log.Println(len(p))
				switch k := x.(type) {
				case primitive.M:
					for i, j := range k {
						projectMap[i] = fmt.Sprint(j)
					}
				}
			}
		}

		for x, y := range projectMap {
			fmt.Println(x, y)
			if x == "project_name" {
				project.ProjectName = y
			}
			if x == "tools_use_as" {
				project.ToolsUseAs = y
			}
			if x == "project_content" {
				project.ProjectContent = y
			}

		}

		c.HTML(http.StatusOK, "show-project.html", gin.H{
			"projectID":      project.ID,
			"projectName":    project.ProjectName,
			"projectContent": project.ProjectContent,
			"toolsUseAs":     project.ToolsUseAs,
		})
	}
}

/*
ModifyUserProject - this method helps to post modified and changes in user projects
and also update the project status as well
*/
func (ts *TrackSpace) ModifyUserProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		sourceLink := c.Param("src")

		projectID := c.Param("id")
		ok := primitive.IsValidObjectID(projectID)
		if sourceLink != "show-project" && !ok {
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: errors.New("invalid ID cannot convert the Object ID")})
		}
		project.ID = projectID
		project.ProjectName = strings.ToLower(c.PostForm("project-name"))
		project.ToolsUseAs = strings.ToLower(c.PostForm("project-tool-use"))
		project.ProjectContent = c.Request.Form.Get("myText")
		project.Status = "modified"
		project.UpdatedAt = time.Now().Format("2006-01-02")
		project.CreatedAt = time.Now().Format("2006-01-02")

		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)

		log.Println(projectID)

		err := ts.tsDB.ModifyProjectData(userData.UserID, projectID, project)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}

		c.HTML(http.StatusSeeOther, "dash.html", gin.H{
			"updateProject": fmt.Sprintf("%s project updated successfully", project.ProjectName),
		})
	}
}

/*
DeleteProject : this is to delete select project existing in the database
*/
func (ts *TrackSpace) DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		sourceLink := c.Param("src")
		ok := primitive.IsValidObjectID(c.Param("id"))
		if sourceLink != "project-table" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		project.ID = c.Param("id")
		ok = primitive.IsValidObjectID(project.ID)
		if !ok {
			log.Println("invalid ID cannot convert the Object ID")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid ID cannot convert the Object ID")})
		}
		err := ts.tsDB.DeleteUserProject(project.ID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}

		c.HTML(http.StatusSeeOther, "dash.html", gin.H{
			"deleteProject": fmt.Sprintf("%s project delete successfully. Go back to dashboard", project.ProjectName),
		})
	}
}

/*
GetTodo - this will help the user to get the todo-page
to set up a schedule
*/
func (ts *TrackSpace) GetTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "todo.html", gin.H{})
	}
}

/*
PostTodoData : this get the user schedule details from the form and store in the
database
*/
func (ts *TrackSpace) PostTodoData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo model.Todo
		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)

		if err := c.Request.ParseForm(); err != nil {
			log.Println("cannot parse the daily task form")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		userID := userData.UserID
		todo.ID = primitive.NewObjectID().Hex()
		todo.ToDoTask = c.Request.Form.Get("task")
		todo.DateSchedule = c.Request.Form.Get("schedule-date")
		todo.StartTime = c.Request.Form.Get("start-time")
		todo.EndTime = c.Request.Form.Get("end-time")
		todo.Status = "Not done"
		// Server side validation of the user input from a form
		if err := ts.AppConfig.Validator.Struct(&todo); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return
			}
		}

		err := ts.tsDB.StoreTodoData(todo, userID)
		if err != nil {
			log.Println("error while inserting todo data in database")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}

		c.HTML(http.StatusOK, "todo.html", gin.H{
			"addTodo": fmt.Sprintf("%s added to schedule plans", todo.ToDoTask),
		})
	}
}

// ShowTodoTable : this  handler direct the user to a page to make changes and modify their
// existing todo store in the database
func (ts *TrackSpace) ShowTodoTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		todo := make(map[string]interface{})
		var allTodo []map[string]interface{}

		tsData := sessions.Default(c)
		userData := tsData.Get("session_data").(model.SessionData)
		userID := userData.UserID

		user, err := ts.tsDB.SendUserDetails(userID)
		if err != nil {
			log.Println("cannot get user project data from the database")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}
		switch p := user["todo"].(type) {
		case primitive.A:
			for _, x := range p {
				switch k := x.(type) {
				case primitive.M:
					for i, j := range k {
						todo[i] = j
					}
					allTodo = append(allTodo, k)
				}
			}
		}
		c.HTML(http.StatusOK, "todo-table.html", gin.H{
			"Todos":     allTodo,
			"FirstName": user["first_name"],
			"LastName":  user["last_name"],
		})
	}
}

/*
ShowTodoSchedule : this will show the selected schedule plans to show all it fulls
details  and as well make changes to it
*/
func (ts *TrackSpace) ShowTodoSchedule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo model.Todo
		TodoMap := make(map[string]string)
		sourceLink := c.Param("src")
		todo.ID = c.Param("id")
		ok := primitive.IsValidObjectID(c.Param("id"))

		if sourceLink != "todo-table" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		TodoData, err := ts.tsDB.GetTodoData(todo.ID)
		if err != nil {
			log.Println("cannot get user project data from the database")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}

		switch p := TodoData["todo"].(type) {
		case primitive.A:
			for _, x := range p {
				log.Println(len(p))
				switch k := x.(type) {
				case primitive.M:
					for i, j := range k {
						TodoMap[i] = fmt.Sprint(j)
					}
				}
			}
		}

		for x, y := range TodoMap {
			// fmt.Println(x, y)
			if x == "to_do_task" {
				todo.ToDoTask = y
			}
			if x == "schedule_date" {
				todo.DateSchedule = y
			}
			if x == "start_time" {
				todo.StartTime = y
			}
			if x == "end_time" {
				todo.EndTime = y
			}
		}

		c.HTML(http.StatusOK, "show-todo.html", gin.H{
			"TodoID":       todo.ID,
			"Task":         todo.ToDoTask,
			"DateSchedule": todo.DateSchedule,
			"StartTime":    todo.StartTime,
			"EndTime":      todo.EndTime,
			"Status":       "Done",
		})
	}
}

/*
ModifyUserTodo - this method helps to post modified and changes in user projects
and also update the todo status as well
*/
func (ts *TrackSpace) ModifyUserTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo model.Todo
		sourceLink := c.Param("src")
		ok := primitive.IsValidObjectID(c.Param("id"))
		if sourceLink != "show-todo" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		todo.ID = c.Param("id")
		ok = primitive.IsValidObjectID(todo.ID)
		if !ok {
			log.Println("invalid ID cannot convert the Object ID")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: errors.New("invalid ID cannot convert the Object ID")})
		}
		todo.ToDoTask = c.Request.Form.Get("task")
		todo.DateSchedule = c.Request.Form.Get("schedule-date")
		todo.StartTime = c.Request.Form.Get("start-time")
		todo.EndTime = c.Request.Form.Get("end-time")
		todo.Status = "Done"

		err := ts.tsDB.ModifyTodoData(todo.ID, todo)
		if err != nil {
			log.Println("Error while storing using user project data")
			return
		}
		c.HTML(http.StatusSeeOther, "dash.html", gin.H{
			"updateTodo": fmt.Sprintf("%s planned schedule changed", todo.ToDoTask),
		})
	}
}

/*
DeleteTodo : this is to delete selected todo task existing in the database
*/
func (ts *TrackSpace) DeleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var todo model.Todo
		sourceLink := c.Param("src")
		ok := primitive.IsValidObjectID(c.Param("id"))
		if sourceLink != "todo-table" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		todo.ID = c.Param("id")
		ok = primitive.IsValidObjectID(todo.ID)
		if !ok {
			log.Println("invalid ID cannot convert the Object ID")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: errors.New("project id is invalid")})
		}
		err := ts.tsDB.DeleteUserTodo(todo.ID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}
		c.HTML(http.StatusSeeOther, "dash.html", gin.H{
			"deleteTodo": fmt.Sprintf(" %s schedule plan delete successfully. Go back to dashboard", todo.ID),
		})
	}
}

// ExecuteLogOut - to log out user from the dashboard
func (ts *TrackSpace) ExecuteLogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		newToken := tsData.Get("refreshToken")
		tsData.Set("newToken", newToken)
		tsData.Clear()
		tsData.Options(sessions.Options{MaxAge: -1})
		_ = tsData.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

/*
AdminPage : this is the Track-space admin page to have a full view of the register user and their
important information as well
*/
func (ts *TrackSpace) AdminPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			TotalProject int
			TotalUser    int
			TotalTodo    int
		)

		tsDoc := make(map[string]interface{})
		var tsUser []map[string]interface{}
		var countryList []string

		documents, err := ts.tsDB.GetAllUserData()
		if err != nil {
			log.Println("cannot get user project data from the database")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}

		TotalUser = len(documents)

		for _, document := range documents {
			for k, v := range document {
				tsDoc[k] = v
				tsDoc["del"] = "delete"
				if k == "project_details" {
					switch p := v.(type) {
					case primitive.A:
						TotalProject = len(p)
					}
				}
				if k == "todo" {
					switch t := v.(type) {
					case primitive.A:
						TotalProject = len(t)
					}
				}
				if k == "country" {
					for _, c := range countryList {
						switch c {
						case fmt.Sprint(v):
							fallthrough

						default:
							countryList = append(countryList, fmt.Sprint(v))
						}
					}
				}

			}

			tsUser = append(tsUser, tsDoc)
		}

		type userStat struct {
			total string
			todo  string
			users string
		}

		stat := userStat{
			total: strconv.Itoa(TotalProject),
			todo:  strconv.Itoa(TotalTodo),
			users: strconv.Itoa(TotalUser),
		}

		statData, err := json.MarshalIndent(stat, "", "  ")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}

		_ = ioutil.WriteFile("./static/json/stat.json", statData, 0o644)

		c.HTML(http.StatusOK, "admin.html", gin.H{
			"tsAdmin": tsUser,
		})
	}
}

/*
AdminDeleteUser : this will allow the admin to delete any user straightaway from the database
and also notify the admin team on the changes made
*/
func (ts *TrackSpace) AdminDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("id")
		src := c.Param("src")
		ok := primitive.IsValidObjectID(userId)
		if src != "admin" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}
		err := ts.tsDB.AdminDeleteUserData(userId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
		}
		message := fmt.Sprintf(`
			<strong>Confirmation for Deleted Account </strong><br>
			Hi, %s:<br>
			<p>
               This is to confirm that you have delete user 
               with an ID: %s from track-space.We hope you 
               are well informed that this action with clear 
               all the user data store on track space
			</p>
			`, "Admin", userId)
		mailMsg := model.Email{
			Subject:  "Confirmation for Deleted Account",
			Content:  message,
			Sender:   "official.trackspace@gmail.com",
			Receiver: "official.trackspace@gmail.com",
			Template: "email.html",
		}

		ts.AppConfig.MailChan <- mailMsg

		c.Redirect(http.StatusSeeOther, "/auth/admin")
	}
}

/*
ChatRoom : this load up the page for user to experience a simple real time communication
*/
func (ts *TrackSpace) ChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	}
}

/*
ChatRoomEndpoint : this is a page for online registered users to have full experience of real time
communication with other user to communicate , interact as well discuss  and share ideas among themselves
share
*/
func (ts *TrackSpace) ChatRoomEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var res wsmodel.SocketResponse
		wsConn, err := ws.UpgradeSocketConn.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Fatal("Unable to connect to socket")
			return
		}

		res.Message = `<p>chatroom handler upgrade</p>`

		connect := wsconfig.SocketConnection{Conn: wsConn}
		ws.Client[connect] = ""

		err = wsConn.WriteJSON(res.Message)
		if err != nil {
			return
		}
		go ws.SendDataToChannel(&connect)
	}
}
