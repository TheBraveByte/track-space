package model

import (
	"time"

	"github.com/gorilla/websocket"
)

type Auth struct {
	Token string
}

type User struct {
	ID             string      `json:"_id" bson:"_id" Usage:"required,alphanumeric"`
	FirstName      string      `json:"first_name" Usage:"required,alpha"`
	LastName       string      `json:"last_name" Usage:"required,alpha"`
	Email          string      `json:"email" Usage:"required,email"`
	Password       string      `json:"password" Usage:"min=8,max=20"`
	YrsOfExp       string      `json:"yrs_of_exp" Usage:"numeric"`
	Country        string      `json:"country" Usage:"required,alpha"`
	PhoneNumber    string      `json:"phone_number" Usage:"required"`
	IPAddress      string      `json:"ip_address"`
	Address        string      `json:"address" Usage:"required"`
	UserType       []string    `son:"user_type"`
	Stack          []string    `json:"stack"`
	ProjectDetails []Project   `json:"project_details" bson:"project_details"`
	Todo           []Todo `json:"todo" bson:"todo"`
	Data           []Data      `json:"data" bson:"data"`
	CreatedAt      time.Time   `json:"created_at" Usage:"datetime=2006-01-02"`
	UpdatedAt      time.Time   `json:"updated_at" Usage:"datetime=2006-01-02"`
	Token          string      `json:"token" Usage:"jwt"`
	RenewToken     string      `json:"renew_token" Usage:"jwt"`
}

type Project struct {
	ID             string `bson:"_id"`
	ProjectName    string `json:"project_name" Usage:"required"`
	ProjectContent string `json:"project_content"`
	ToolsUseAs     string `json:"tools_use_as" Usage:"required"`
	UpdatedAt      string `json:"updated_at"`
	CreatedAt      string `json:"created_at"`
	Status         string `json:"status"`
}

type Data struct {
	ID      string `json:"_id" bson:"_id"`
	Date    string `json:"date"`
	Code    int    `json:"code"`
	Article int    `json:"article"`
	Text    int    `json:"text"`
	Todo    int    `json:"todo"`
	Total   int    `json:"total"`
}

type Email struct {
	ID       string `json:"_id" bson:"_id"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Receiver string `json:"receiver" Usage:"required"`
	Sender   string `json:"sender" Usage:"required"`
	Template string `json:"template"`

}

type Todo struct {
	ID           string `json:"_id" bson:"_id"`
	ToDoTask     string `json:"to_do_task"`
	DateSchedule string `json:"date_schedule"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Status         string `json:"status"`

}

/*Working with Web Socket*/

var (
	WebSkChan = make(chan SocketPayLoad)
	Client    = make(map[SocketConnection]string)
)

type SocketConnection struct {
	*websocket.Conn
}

type SocketPayLoad struct {
	Condition  string           `json:"condition"`
	Message    string           `json:"message"`
	UserName   string           `json:"user_name"`
	SocketConn SocketConnection `json:"-"`
}

type SocketResponse struct {
	Condition   string `json:"condition"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
	ConnectedUSer []string `json:"connected_user"`
}
