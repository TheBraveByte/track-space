package model

import (
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Auth struct {
	Token string
}

// type User struct {
// 	ID             primitive.ObjectID `json:"_id" bson:"_id"`
// 	FirstName      string             `form:"first_name" json:"first_name" Usage:"required max=32 min=3" binding:"required max=32 min=3"`
// 	LastName       string             `form:"last_name" json:"last_name" Usage:"required max=32 min=3" binding:"required max=32 min=3"`
// 	Email          string             `form:"email" json:"email" Usage:"email required" binding:"required max=32 min=3" `
// 	Password       string             `form:"password" json:"password" Usage:"required_with=Email alphanum" binding:"required max=32 min=3"`
// 	YrsOfExp       string             `form:"yrs_of_exp" json:"yrs_of_exp" Usage:"numeric omitempty"`
// 	Country        string             `form:"country" json:"country" Usage:"required" binding:"required max=32 min=3"`
// 	PhoneNumber    string             `form:"phone_number" json:"phone_number" Usage:"required max=15 min=8" binding:"required max=32 min=3"`
// 	IPAddress      string             `form:"ip_address" json:"ip_address"`
// 	Address        string             `form:"address" json:"address" Usage:"required" binding:"required max=32 min=3"`
// 	UserType       []string           `form:"user_type" json:"user_type" Usage:"omitempty"`
// 	Stack          []string           `form:"stack" json:"stack" Usage:"omitempty"`
// 	ProjectDetails []Project          `form:"project_details" json:"project_details" bson:"project_details"`
// 	CreatedAt      time.Time          `json:"created_at"`
// 	UpdatedAt      time.Time          `json:"updated_at"`
// 	Token          string             `json:"token"`
// 	RenewToken     string             `json:"renew_token"`
// }

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	YrsOfExp       string             `json:"yrs_of_exp"`
	Country        string             `json:"country"`
	PhoneNumber    string             `json:"phone_number"`
	IPAddress      string             `json:"ip_address"`
	Address        string             `json:"address"`
	UserType       []string           `son:"user_type"`
	Stack          []string           `son:"stack"`
	ProjectDetails []Project          `json:"project_details" bson:"project_details"`
	Todo           []DailyTask        `json:"todo" bson:"todo"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Token          string             `json:"token"`
	RenewToken     string             `json:"renew_token"`
}

type Project struct {
	ID             primitive.ObjectID `bson:"_id"`
	ProjectName    string             `json:"project_name"`
	ProjectContent string             `json:"project_content"`
	ToolsUseAs     string             `json:"tools_use_as"`
	StartTime      time.Time          `json:"start_time"`
	EndTime        time.Time          `json:"end_time"`
	Duration       time.Duration      `json:"duration"`
}

type Email struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Message      string             `json:"message"`
	Receiver     string             `json:"receiver" validate:"required"`
	Sender       string             `json:"sender" validate:"required"`
	MailTemplate string             `json:"mail_template"`
}

type DailyTask struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	ToDoTask     string             `json:"to_do_task"`
	DateSchedule string             `json:"date_schedule"`
	StartTime    string             `json:"start_time"`
	EndTime      string             `json:"end_time"`
}



/*Working with Web Socket*/

var WebSkChan = make(chan SocketPayLoad)
var Client = make(map[SocketConnection]string)


type SocketConnection struct{
	*websocket.Conn
}

type SocketPayLoad struct{
	Condition string `json:"condition"`
	Message string `json:"message"`
	UserName string `json:"user_name"`
	SocketConn SocketConnection `json:"-"`

	
}

type SocketResponse struct{
	Condition string `json:"condition"`
	Message string `json:"message"`
	MessageType string `json:"message_type"`
	UserName string `json:"user_name"`
	ConnectedUSer []string `json:"connected_user"`

}