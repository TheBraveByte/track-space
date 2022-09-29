package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yusuf/track-space/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app config.AppConfig

func TestTrackSpace_Article(t *testing.T) {

	var tests = []struct {
		name         string
		appConfig    *config.AppConfig
		count        map[string]int
		countArticle int
	}{
		{"Article", nil, map[string]int{"article": 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTestTrackSpace(tt.appConfig)
			ts.Article(tt.count, tt.countArticle)
		})
	}
}

func TestTrackSpace_Code(t *testing.T) {

	var tests = []struct {
		name      string
		appConfig *config.AppConfig
		count     map[string]int
		countCode int
	}{
		{"code", nil, map[string]int{"code": 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTestTrackSpace(tt.appConfig)
			ts.Code(tt.count, tt.countCode)
		})
	}
}
func TestTrackSpace_Text(t *testing.T) {

	var tests = []struct {
		name      string
		appConfig *config.AppConfig
		count     map[string]int
		countText int
	}{
		{"text", nil, map[string]int{"text": 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTestTrackSpace(tt.appConfig)
			ts.Text(tt.count, tt.countText)
		})
	}
}

func TestTrackSpace_Todo(t *testing.T) {

	var tests = []struct {
		name      string
		appConfig *config.AppConfig
		count     map[string]int
		countTodo int
	}{
		{"todo", nil, map[string]int{"todo": 2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTestTrackSpace(tt.appConfig)
			ts.Todo(tt.count, tt.countTodo)
		})
	}
}

func TestTrackSpace_HomePage(t *testing.T) {

	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"Homepage", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/", ts.HomePage(), func(context *gin.Context) {
				context.HTML(http.StatusOK, "home-page.html", gin.H{"authenticate": 0})
			})
			rq, _ := http.NewRequest("GET", "/", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestTrackSpace_Contact(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"Contact", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/contact", ts.Contact(), func(context *gin.Context) {
				context.HTML(http.StatusOK, "contact.html", gin.H{})
			})
			rq, _ := http.NewRequest("GET", "/contact", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

//func TestTrackSpace_PostContact(t *testing.T) {
//	tests := []struct {
//		name       string
//		val        url.Values
//		AppConfig  *config.AppConfig
//		statusCode int
//	}{
//		{"post_contact", url.Values{
//			"fullName": []string{"track-space"},
//			"email":    []string{"trackspace@admin.com"},
//			"msg":      []string{"Hello world"},
//		}, &app, http.StatusOK},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			w := httptest.NewRecorder()
//			router := TrackSpaceSetUp()
//			ts := NewTestTrackSpace(tt.AppConfig)
//			router.POST("/contact", ts.PostContact(), func(context *gin.Context) {
//				if err := context.Request.ParseForm(); err != nil {
//					log.Panic("form not parsed")
//					return
//				}
//				name := tt.val.Get("fullName")
//				email := tt.val.Get("email")
//				msg := tt.val.Get("msg")
//				if name != "" && email != "" && msg != "" {
//					context.HTML(http.StatusOK, "contact.html", gin.H{})
//				}
//			})
//			rq, _ := http.NewRequest("POST", "/contact", strings.NewReader(tt.val.Encode()))
//			rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//			router.ServeHTTP(w, rq)
//			assert.Equal(t, tt.statusCode, w.Code)
//		})
//	}
//}

func TestTrackSpace_SignUpPage(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"sign-up", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/signup", ts.SignUpPage())
			rq, _ := http.NewRequest("GET", "/signup", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestTrackSpace_GetUserInfo(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"user-info", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/user-info", ts.GetUserInfo())
			rq, _ := http.NewRequest("GET", "/user-info", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestTrackSpace_GetLoginPage(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"login", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/login", ts.GetLoginPage())
			rq, _ := http.NewRequest("GET", "/login", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
func TestTrackSpace_ResetPassword(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"reset password", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/reset-password", ts.ResetPassword())
			rq, _ := http.NewRequest("GET", "/reset-password", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestTrackSpace_ProjectWorkspace(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"workspace", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/auth/user/workspace", ts.ProjectWorkspace())
			rq, _ := http.NewRequest("GET", "/auth/user/workspace", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
func TestTrackSpace_GetTodo(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"todo", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/auth/user/todo", ts.GetTodo())
			rq, _ := http.NewRequest("GET", "/auth/user/todo", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestTrackSpace_ExecuteLogOut(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"log-out", &app, http.StatusTemporaryRedirect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/user/log-out", ts.ExecuteLogOut())
			rq, _ := http.NewRequest("GET", "/user/log-out", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
func TestTrackSpace_ChatRoom(t *testing.T) {
	tests := []struct {
		name       string
		AppConfig  *config.AppConfig
		statusCode int
	}{
		{"chat-room", &app, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router := TrackSpaceSetUp()
			ts := NewTestTrackSpace(tt.AppConfig)
			router.GET("/user/chat", ts.ChatRoom())
			rq, _ := http.NewRequest("GET", "/user/chat", nil)
			router.ServeHTTP(w, rq)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
