package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/yusuf/track-space/pkg/data"
	"github.com/yusuf/track-space/pkg/data/tsRepoStore"
	"github.com/yusuf/track-space/pkg/key"
	"github.com/yusuf/track-space/pkg/temp"
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

func (ts *TrackSpace) HomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var templateData temp.TemplateData
		c.HTML(http.StatusOK, "home-page.html", gin.H{
			"authenticate": 0,
		})
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
		user.ID = primitive.NewObjectID().Hex()
		user.Email = c.Request.Form.Get("email")
		user.Password = key.HashPassword(c.Request.Form.Get("password"))
		log.Println(user.ID, user.Email, user.Password)

		//Server side validation of the user input from a form
		if err := Validate.Struct(user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok{
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return 
			}
		}

		tsData := sessions.Default(c)
		tsData.Set("email", user.Email)
		tsData.Set("password", user.Password)
		tsData.Set("userID", user.ID)

		count, err := ts.tsDB.InsertUserInfo(user.ID, user.Email, user.Password)
		if err != nil {
			log.Println(err)
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		if err = tsData.Save(); err != nil {
			log.Println("error from the session storage")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
			return
		}

		if count == 1 {
			c.HTML(http.StatusOK, "login-page.html", gin.H{
				"msg": "You have previously sign-up\nlog-in into your account",
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
		userID := fmt.Sprint(tsData.Get("userID"))

		if err := c.Request.ParseForm(); err != nil {
			log.Println(err)
			return
		}
		user.FirstName = c.Request.Form.Get("first-name")
		user.LastName = c.Request.Form.Get("last-name")
		user.Address = c.Request.Form.Get("address")
		user.YrsOfExp = c.Request.Form.Get("yrs-of-exp")
		user.Country = c.Request.Form.Get("nation")
		user.PhoneNumber = c.Request.Form.Get("phone")
		user.Stack = append(user.Stack, c.Request.Form.Get("stack-name"))
		user.IPAddress = c.Request.RemoteAddr
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		t1, t2 := "", ""

		//Server side validation of the user input from a form
		if err := Validate.Struct(user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok{
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return 
			}
		}
		err := ts.tsDB.UpdateUserInfo(user, userID, t1, t2)
		if err != nil {
			log.Println("Cannot update user info")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}
		message := fmt.Sprintf(`
			<strong>Confirmation for Account Created</strong><br>
			Hi, %s:<br>
			<p>This is to confirm that your have sign up for track-space.
			We hope you have a wonderful experience using our platform as
			your workspace for your project
			</p>
			`,user.FirstName)
		mailMsg := model.Email{
			Subject: "Confirmation for Account Created",
			Content: message,
			Sender:"yusufakinleye@gmail.com",
			Receiver: user.Email,
			Template: "email.html",
		}

		ts.AppConfig.MailChan <- mailMsg

		TeamMessage := fmt.Sprintf(`
			<strong>New User Account Notification</strong><br>
			Hi, %s:<br>
			<p>This is notify you guys that new user with an <strong> ID:</strong> %s 
			and <strong>IPAddress : </strong> of %s sign up for track-space.
			</p>
			`,"track-space Team",userID, user.IPAddress)
		TeamMailMsg := model.Email{
			Subject: "Confirmation for Account Created",
			Content: TeamMessage,
			Sender:"yusufakinleye@gmail.com",
			Receiver: "yusufakinleye@gmail.com",
			Template: "email.html",
		}

		ts.AppConfig.MailChan <- TeamMailMsg
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
		var templateData temp.TemplateData
		templateData.IsAuthenticated = 1

		if err := c.Request.ParseForm(); err != nil {
			log.Println("error while parsing form")
			return
		}

		email := fmt.Sprint(tsData.Get("email"))
		password := fmt.Sprint(tsData.Get("password"))
		userID := fmt.Sprint(tsData.Get("userID"))
		fmt.Println(userID)
		IPAddress := c.Request.RemoteAddr

		// Posted form value
		var user model.User
		user.Email = c.Request.Form.Get("email")
		user.Password = c.Request.Form.Get("password")

		//Server side validation of the user input from a form
		if err := Validate.Struct(user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok{
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return 
			}
		}

		// check to verify for the stored hashed password in database
		if email == user.Email && password == user.Password {
			ok, _ := ts.tsDB.VerifyLogin(userID, password, user.Password)
			if ok {
				// check to match hashed password and the user password input
				token, newToken, err := auth.GenerateJWTToken(user.Email, userID, IPAddress)
				if err != nil {
					log.Println("cannot generate json web token")
					_ = c.AbortWithError(http.StatusBadRequest, gin.Error{
						Err:  err,
						Type: 0,
						Meta: nil,
					})
					return
				}

				authData := templateData.AuthData
				authData["auth"] = []string{token, newToken}
				tokenGen := authData["auth"][0]
				newTokenGen := authData["auth"][1]
				log.Println(tokenGen)

				err = ts.tsDB.UpdateUserField(userID, tokenGen, newTokenGen)
				if err != nil {
					log.Println("cannot update user info")
					return
				}
				c.SetCookie("bearerToken", tokenGen, 60*60*24*1200, "/", "localhost", false, true)
				log.Println("Successfully login")

			}
		} else {
			c.HTML(http.StatusOK, "home-page.html", gin.H{
				"error": "incorrect email or password",
			})
		}
		if err := tsData.Save() ; err != nil {
			log.Println("error from the session storage")
		}
		c.HTML(http.StatusOK, "home-page.html", gin.H{
			"success":      "successfully login! click dashboard",
			"authenticate": templateData.IsAuthenticated,
		})
	}
}

// Todo Each of the  functions are use to organize the user data in the database
// Todo
func (ts *TrackSpace) Todo(count map[string]int, countTodo int) {
	count["todoNo"] = countTodo
}

func (ts *TrackSpace) Code(count map[string]int, countText int) {
	count["text"] = countText
}

func (ts *TrackSpace) Text(count map[string]int, countCode int) {
	count["code"] = countCode
}

func (ts *TrackSpace) Article(count map[string]int, countArticle int) {
	count["article"] = countArticle
}

// GetDashBoard - this show the user dashboard with respect to all the database
// details and queries; full brief or user activities
func (ts *TrackSpace) GetDashBoard() gin.HandlerFunc {
	// a lot of logic will be done here ..... a lot
	return func(c *gin.Context) {
		t, ok := c.Get("token")
		if ok {
			tsData := sessions.Default(c)
			userID := fmt.Sprint(tsData.Get("userID"))
			log.Println(userID)

			user, err := ts.tsDB.SendUserDetails(userID)
			if err != nil {
				log.Println(err)
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
										if i == "created_at"{
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

			code := count["code"]
			text := count["text"]
			article := count["article"]
			todo := count["todoNo"]
			totalProjects := count["article"] + count["text"] + count["code"]

			var data model.Data
			data.Article = article
			data.Code = code
			data.Date = currentDate
			data.Text = text
			data.Todo = todo
			data.Total = totalProjects

			if currentDate == storedDate{
				_ = ts.tsDB.UpdateUserStat(data, userID)

			}

			// this controller with still be updated as I progress
			if err := tsData.Save(); err != nil {
				log.Println("error from the session storage")
			}
			c.HTML(http.StatusOK, "dash.html", gin.H{
				"FirstName": user["first_name"],
				"LastName":  user["last_name"],
				"token":     t,
			})
		}
	}
}

// WorkSpace -  this show the user workspace/ worksheet to execute
// projects and also to make use of other tools
func (ts *TrackSpace) WorkSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "work.html", gin.H{})
	}
}

// PostWorkSpace - this will validate the project model and help to insert
// the projects details in the database
func (ts *TrackSpace) PostWorkSpace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		tsData := sessions.Default(c)
		userID := fmt.Sprint(tsData.Get("userID"))

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

		//Server side validation of the user input from a form
		if err := Validate.Struct(&project); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok{
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return 
			}
		}

		err := ts.tsDB.StoreProjectData(userID, project)
		if err != nil {
			log.Println("Error while storing using user project data")
			return
		}

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}
		c.HTML(http.StatusOK, "work.html", gin.H{
			"save": fmt.Sprintf("%v suubmitted successfully", project.ProjectName),
		})
	}
}

/* 
  DailyTaskTodo - this will help the user to get the todo-page
  to set up a schedule
*/
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

		if err := c.Request.ParseForm(); err != nil {
			log.Println("cannot parse the daily task form")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		userID := fmt.Sprintf("%s", tsData.Get("userID"))
		task.ID = primitive.NewObjectID().Hex()
		task.ToDoTask = c.Request.Form.Get("task")
		task.DateSchedule = c.Request.Form.Get("date_schedule")
		task.StartTime = c.Request.Form.Get("start-time")
		task.EndTime = c.Request.Form.Get("end-time")

		//Server side validation of the user input from a form
		if err := Validate.Struct(&task); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok{
				_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				log.Println(err)
				return 
			}
		}

		err := ts.tsDB.StoreDailyTaskData(task, userID)
		if err != nil {
			log.Println("error while inserting todo data in database")
			_ = c.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		if err := tsData.Save(); err != nil {
			log.Println("error from the session storage")
		}

		c.HTML(http.StatusOK, "daily-task.html", gin.H{
			"taskSaved": "Schedule task added",
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
		userID := fmt.Sprintf("%s", tsData.Get("userID"))

		fmt.Println(tsData.Get("userID"))
		fmt.Println(userID)

		user, err := ts.tsDB.SendUserDetails(userID)
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
			"Status":    user["status"],
		})
	}
}

// ShowUserProject : this  handler direct the user to a page to make changes and modify their
// existing projects store in the database
func (ts *TrackSpace) ShowUserProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		projectMap := make(map[string]string)
		sourceLink := c.Param("src")
		projectID := c.Param("id")
		// projectName := c.Param("project-name")
		ok := primitive.IsValidObjectID(c.Param("id"))

		if sourceLink != "project-table" && !ok{
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		projectData, err := ts.tsDB.GetProjectData(projectID)
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
			fmt.Println(x,y)
			if x == "project_name" {
				project.ProjectName = y
			} 
			if x == "tools_use_as"{
				project.ToolsUseAs = y
			}
			if x == "project_content"{
				project.ProjectContent = y
			}
			if x == "_id"{
				project.ID = y
			}
		}

		c.HTML(http.StatusOK, "show-project.html", gin.H{
			"projectID" : project.ID,
			"projectName" : project.ProjectName,
			"projectContent": project.ProjectContent,
			"toolsUseAs" : project.ToolsUseAs,

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
		ok := primitive.IsValidObjectID(c.Param("id"))
		if sourceLink != "show-project" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		projectID := c.Param("id")
		ok = primitive.IsValidObjectID(projectID)
		if !ok {
			log.Println("invalid ID cannot convert the Object ID")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: errors.New("project id is invalid")})
		}
		project.ProjectName = strings.ToLower(c.PostForm("project-name"))
		project.ToolsUseAs = strings.ToLower(c.PostForm("project-tool-use"))
		project.ProjectContent = c.Request.Form.Get("myText")
		project.Status = "modified"

		err := ts.tsDB.ModifyProjectData(projectID, project)
		if err != nil {
			log.Println("Error while storing using user project data")
			return
		}
		c.HTML(http.StatusOK, "show-project.html", gin.H{})
	}
}


func(ts *TrackSpace) DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var project model.Project
		sourceLink := c.Param("src")
		ok := primitive.IsValidObjectID(c.Param("id"))
		if sourceLink != "project" && !ok {
			_ = c.AbortWithError(http.StatusInternalServerError, gin.Error{Err: errors.New("invalid url parameters")})
			return
		}

		project.ID = c.Param("id")
		ok = primitive.IsValidObjectID(project.ID)
		if !ok {
			log.Println("invalid ID cannot convert the Object ID")
			_ = c.AbortWithError(http.StatusNotFound, gin.Error{Err: errors.New("project id is invalid")})
		}


	}
}
// SettingPage - handlers to make general changes to user platform
// dashboard
func (ts *TrackSpace) SettingPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// is to together settings html templates
		c.HTML(http.StatusOK, "setting.html", gin.H{})
	}
}

// PostSettingChange - to execute and implemented the change in  settings
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
		err := tsData.Save()
		if err != nil {
			return
		}
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (ts *TrackSpace) ChatRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", gin.H{})
	}
}

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

		err = wsConn.WriteJSON(res.Message)
		if err != nil {
			return
		}
	}
}
