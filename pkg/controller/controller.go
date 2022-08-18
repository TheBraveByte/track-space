package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/yusuf/track-space/pkg/data"
	"github.com/yusuf/track-space/pkg/data/tsRepoStore"
	"github.com/yusuf/track-space/pkg/key"
	"github.com/yusuf/track-space/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/yusuf/track-space/pkg/auth"
	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/model"
)

// Validate - to help check for a validated json database model
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

// SignUpPage - Handler to get the sign up page for user
func (ts *TrackSpace) SignUpPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup-page.html", gin.H{})
	}
}

// PostSignUpPage - this validate the user input and store the value in
// session as cookies for future usage, insert user input in the database
// and also check for existing user
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

// GetUserInfo - handler to get the user-details/ info page
func (ts *TrackSpace) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-info.html", gin.H{})
	}
}

// PostUserInfo - this validdate the user model and get the user input details
// and store the details in the database and this helps to redirect
// the user to a login page
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

		// info := make(map[string]interface{})
		email := tsData.Get("email")
		user.ID = primitive.NewObjectID()
		user.FirstName = c.Request.Form.Get("first-name")
		user.LastName = c.Request.Form.Get("last-name")
		user.Address = c.Request.Form.Get("address")
		user.YrsOfExp = c.Request.Form.Get("yrs-of-exp")
		user.Country = c.Request.Form.Get("nation")
		user.PhoneNumber = c.Request.Form.Get("phone")
		user.Stack = append(user.Stack, c.Request.Form.Get("stack-name"))

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		t1, t2 := "", ""

		err = ts.tsDB.UpdateUserInfo(user, email, t1, t2)
		if err != nil {
			log.Println("Cannot update user info")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
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

// PostLoginPage : this handler help to verify the user password, authenicate other
// user login details with respect to the database,generate a authorization token
// for the user, as well as authorize the user and set the Response Header
// with the Bearer Token
func (ts *TrackSpace) PostLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		var user model.User

		if err := c.Request.ParseForm(); err != nil {
			log.Println("error while parsing form")
			return
		}

		email := tsData.Get("email")
		password := tsData.Get("password")

		IPAddress := c.Request.RemoteAddr

		// Posted form value
		postEmail := c.Request.Form.Get("email")
		postPassword := c.Request.Form.Get("password")

		// Previous store details in the database
		ok, hashedPassword := ts.tsDB.VerifyLogin(postEmail)

		if ok {
			_, msg := key.VerifyPassword(postPassword, hashedPassword)
			log.Println(msg)
			if postPassword == password && postEmail == email {
				token, newToken, err := auth.GenerateJWTToken(email, password, IPAddress)
				if err != nil {
					log.Println("cannot generate json web token")
					_ = c.AbortWithError(http.StatusBadRequest, gin.Error{
						Err:  err,
						Type: 0,
						Meta: nil,
					})
					return
				}

				authData := make(map[string][]string)
				authData["auth"] = []string{token, newToken}
				t1 := authData["auth"][0]
				t2 := authData["auth"][1]
				fmt.Println(t1, t2)
				err = ts.tsDB.UpdateUserInfo(user, email, t1, t2)
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
		fmt.Println("Successfully login")

		c.HTML(http.StatusOK, "/", gin.H{})
	}
}

// GetDashBoard - this show the user dashboard with respect to all the database
// details and queries; full brief or user activites
func (ts *TrackSpace) GetDashBoard() gin.HandlerFunc {
	// a lot of logic will be done here ..... a lot
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		email := fmt.Sprintf("%s", tsData.Get("email"))

		user, err := ts.tsDB.SendUserDetails(email)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}

		// this controller with still be updated as I progress

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

// WorkSpace -  this show the user workspace/ worksheet to execute
// projects and also to make use of other tools
func (ts *TrackSpace) WorkSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "work.html", gin.H{
			"StartTime": StartTime.Local().UTC(),
		})
	}
}

// PostWorkSpace - this will validate the project model and help to insert
// the projects details in the database
func (ts *TrackSpace) PostWorkSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		tsData := sessions.Default(c)
		userEmail := tsData.Get("email")

		if validateErr := Validate.Struct(&project); validateErr != nil {
			log.Println("cannot validate project struct")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: validateErr})
			return
		}

		// Getting the project data
		if err := c.Request.ParseForm(); err != nil {
			log.Println("cannot parse the workspace form")
			return
		}
		project.ProjectName = strings.ToTitle(c.PostForm("project-name"))
		project.ToolsUseAs = strings.ToLower(c.PostForm("project-tool-use"))
		project.ProjectContent = c.PostForm("editor")
		project.StartTime = StartTime.Local().UTC()

		err := ts.tsDB.StoreWorkSpaceData(userEmail, project)
		if err != nil {
			log.Println("Error while storing using user project data")
			return
		}
		tsData.AddFlash("successfully submitted project")

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}
		c.HTML(http.StatusOK, "work.html", gin.H{
			"Save": tsData.Flashes("successfully submitted project"),
		})
	}
}

// ProcessWorkSpace - this will help execute the queries to help group different
// type of projects in a map[string]interface{}
func (ts *TrackSpace) ProcessWorkSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		var projectData model.User
		if err := c.ShouldBind(&projectData); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if ValidateErr := Validate.Struct(&projectData); ValidateErr != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		projectData.Email = fmt.Sprintf("%s", tsData.Get("email"))

		count, err := ts.tsDB.OrganizeWorkSpaceData(projectData.Email)
		if err != nil {
			log.Println(err)
			return
		}
		code := count["code"]
		text := count["text"]

		c.HTML(http.StatusOK, "dash.html", gin.H{
			"CodeCount": code,
			"TextCount": text,
		})
	}
}

// DailyTaskTodo - this will help the user to get the todo-page
//  to set up a schedule
func (ts *TrackSpace) DailyTaskTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "daily-task.html", gin.H{})
	}
}

// PostDailyTaskTodo - this get the user schedule details from the form and store in the
// database
func (ts *TrackSpace) PostDailyTaskTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var task model.DailyTask
		tsData := sessions.Default(c)

		if validateErr := Validate.Struct(&task); validateErr != nil {
			log.Println("cannot validate daily task struct")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: validateErr})
			return
		}

		if err := c.Request.ParseForm(); err != nil {
			log.Println("cannot parse the daily task form")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		userEmail := fmt.Sprintf("%s", tsData.Get("email"))
		task.ToDoTask = c.Request.Form.Get("task")
		task.DateSchedule = c.Request.Form.Get("date_schedule")
		task.StartTime = c.Request.Form.Get("start-time")
		task.EndTime = c.Request.Form.Get("end-time")

		err := ts.tsDB.StoreDailyTaskData(task, userEmail)
		if err != nil {
			log.Println("error while inserting todo data in database")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}
		tsData.AddFlash("successfully submitted project")

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}

		c.HTML(http.StatusOK, "daily-task.html", gin.H{
			"Save": tsData.Flashes("successfully submitted project"),
		})
	}
}

// ShowProjectTable - this give the full projects details of a particular
// user and help the user to modify each projects
func (ts *TrackSpace) ShowProjectTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		proj := make(map[string]interface{})
		tsData := sessions.Default(c)
		userEmail := fmt.Sprintf("%s", tsData.Get("email"))

		userProject, err := ts.tsDB.SendUserDetails(userEmail)
		if err != nil {
			log.Println("cannot get user project data from the database")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: err})
			return
		}
		switch p := userProject["project_details"].(type) {
		case []model.Project:
			for _, x := range p {
				proj["ProjectContent"] = x.ProjectContent
				proj["ProjectName"] = x.ProjectName
				proj["ToolsUseAs"] = x.ToolsUseAs
				proj["StartTime"] = x.StartTime
				proj["EndTime"] = x.EndTime
				proj["Duration"] = x.Duration
				proj["ID"] = x.ID
			}
		default:
		}

		c.HTML(http.StatusOK, "project-table.html", gin.H{
			"project":   proj,
			"FirstName": userProject["first_name"],
			"LastName":  userProject["last_name"],
		})
	}
}

// ShowUserProject : this  handler direct the user to a page to make changes and modify their
// existing projects store in the database
func (ts *TrackSpace) ShowUserProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceLink := c.Param("src")
		ok := primitive.IsValidObjectID(c.Param("id"))

		if sourceLink != "show-project" && !ok {
			log.Println("Invalid url parameters")
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Meta: "invalid parameter"})
			return
		}

		projectID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			log.Println("invalid ID cannot convert the Object ID")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
		}

		data, err := ts.tsDB.GetProjectData(projectID)
		if err != nil {
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
		}
		c.HTML(http.StatusOK, "show-project.html", gin.H{
			"ProjectContent": data["project_content"],
			"ProjectName":    data["project_name"],
			"ToolsUseAs":     data["ToolsUseAs"],
		})
	}
}

// SettingsPage - handlers to make general changes to user platform
// dashboard
func (ts *TrackSpace) SettingPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// is to getbthe settings htmml templates
		c.HTML(http.StatusOK, "setting.html", gin.H{})
	}
}

// PosteSettingChange - to execute and implemented the change in  settings
// page
func (ts *TrackSpace) PostSettingChange() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Write some settings logic program
	}
}

// ExecuteLogOut - to log out user from the dashboard
func (ts *TrackSpace) ExecuteLogOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		tsData := sessions.Default(c)
		tsData.Clear()
		tsData.Options(sessions.Options{MaxAge: -1})
		tsData.Save()
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (ts *TrackSpace) Statistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Statistic report base on the projects and the ToDo schedules
		// implemeted using D3.js
		var _ model.User
		c.HTML(http.StatusOK, "stat.html", gin.H{})
	}
}

func (ts *TrackSpace) ChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	}
}


//
func (ts *TrackSpace) ChatRoomEndpoint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var res model.SocketResponse
		wsConn, err := ws.UpgradeSocketConn.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Fatal("Unable to connect to socket")
			return
		}
		connect := model.SocketConnection{Conn: wsConn}
		model.Client[connect] = ""

		go ws.SendDataToChannel(&connect, ctx)

		wsConn.WriteJSON(res.Message)
	}
}
