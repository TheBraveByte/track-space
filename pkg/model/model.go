package model

type Auth struct {
	Token string
}

// User : Master struct model for user
type User struct {
	ID             string    `bson:"_id" Usage:"required,alphanumeric"`
	FirstName      string    `bson:"first_name" Usage:"required,alpha"`
	LastName       string    `bson:"last_name" Usage:"required,alpha"`
	Email          string    `bson:"email" Usage:"required,email"`
	Password       string    `bson:"password" Usage:"min=8,max=20"`
	YrsOfExp       string    `bson:"yrs_of_exp" Usage:"numeric"`
	Country        string    `bson:"country" Usage:"required,alpha"`
	PhoneNumber    string    `bson:"phone_number" Usage:"required"`
	IPAddress      string    `bson:"ip_address"`
	Address        string    `bson:"address" Usage:"required"`
	Profession     string    `bson:"profession"`
	Stack          []string  `bson:"stack"`
	ProjectDetails []Project `bson:"project_details"`
	Todo           []Todo    `bson:"todo"`
	Data           []Data    `bson:"data"`
	CreatedAt      string    `bson:"created_at" Usage:"datetime=2006-01-02"`
	UpdatedAt      string    `bson:"updated_at" Usage:"datetime=2006-01-02"`
	Token          string    `bson:"token" Usage:"jwt"`
	RenewToken     string    `bson:"renew_token" Usage:"jwt"`
}

// Project : Struct model for user project
type Project struct {
	ID             string `bson:"_id"`
	ProjectName    string `bson:"project_name" Usage:"required"`
	ProjectContent string `bson:"project_content"`
	ToolsUseAs     string `bson:"tools_use_as" Usage:"required"`
	UpdatedAt      string `bson:"updated_at"`
	CreatedAt      string `bson:"created_at"`
	Status         string `bson:"status"`
}

// Data : Struct model to navigate all user activity
type Data struct {
	ID      string `bson:"_id"`
	Date    string `bson:"date"`
	Code    int    `bson:"code"`
	Article int    `bson:"article"`
	Text    int    `bson:"text"`
	Todo    int    `bson:"todo"`
	Total   int    `bson:"total"`
}

// Email : struct model to transmit mail to user and admin
type Email struct {
	ID       string `bson:"_id"`
	Subject  string `bson:"subject"`
	Content  string `bson:"content"`
	Receiver string `bson:"receiver" Usage:"required"`
	Sender   string `bson:"sender" Usage:"required"`
	Template string `bson:"template"`
}

// Todo : struct model for todo schedule for use
type Todo struct {
	ID           string `bson:"_id"`
	ToDoTask     string `bson:"to_do_task"`
	DateSchedule string `bson:"schedule_date"`
	StartTime    string `bson:"start_time"`
	EndTime      string `bson:"end_time"`
	Status       string `bson:"status"`
}

type SessionData struct {
	UserID   string
	Email    string
	Password string
}
