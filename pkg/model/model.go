package model

type Auth struct {
	Token string
}

// User : Master struct model for user
type User struct {
	ID             string    `json:"_id" bson:"_id" Usage:"required,alphanumeric"`
	FirstName      string    `json:"first_name" Usage:"required,alpha"`
	LastName       string    `json:"last_name" Usage:"required,alpha"`
	Email          string    `json:"email" Usage:"required,email"`
	Password       string    `json:"password" Usage:"min=8,max=20"`
	YrsOfExp       string    `json:"yrs_of_exp" Usage:"numeric"`
	Country        string    `json:"country" Usage:"required,alpha"`
	PhoneNumber    string    `json:"phone_number" Usage:"required"`
	IPAddress      string    `json:"ip_address"`
	Address        string    `json:"address" Usage:"required"`
	Profession     string    `json:"profession"`
	Stack          []string  `json:"stack"`
	ProjectDetails []Project `json:"project_details" bson:"project_details"`
	Todo           []Todo    `json:"todo" bson:"todo"`
	Data           []Data    `json:"data" bson:"data"`
	CreatedAt      string    `json:"created_at" Usage:"datetime=2006-01-02"`
	UpdatedAt      string    `json:"updated_at" Usage:"datetime=2006-01-02"`
	Token          string    `json:"token" Usage:"jwt"`
	RenewToken     string    `json:"renew_token" Usage:"jwt"`
}

// Project : Struct model for user project
type Project struct {
	ID             string `bson:"_id"`
	ProjectName    string `json:"project_name" Usage:"required"`
	ProjectContent string `json:"project_content"`
	ToolsUseAs     string `json:"tools_use_as" Usage:"required"`
	UpdatedAt      string `json:"updated_at"`
	CreatedAt      string `json:"created_at"`
	Status         string `json:"status"`
}

// Data : Struct model to navigate all user activity
type Data struct {
	ID      string `json:"_id" bson:"_id"`
	Date    string `json:"date"`
	Code    int    `json:"code"`
	Article int    `json:"article"`
	Text    int    `json:"text"`
	Todo    int    `json:"todo"`
	Total   int    `json:"total"`
}

// Email : struct model to transmit mail to user and admin
type Email struct {
	ID       string `json:"_id" bson:"_id"`
	Subject  string `json:"subject"`
	Content  string `json:"content"`
	Receiver string `json:"receiver" Usage:"required"`
	Sender   string `json:"sender" Usage:"required"`
	Template string `json:"template"`
}

// Todo : struct model for todo schedule for use
type Todo struct {
	ID           string `json:"_id" bson:"_id"`
	ToDoTask     string `json:"to_do_task"`
	DateSchedule string `json:"schedule_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Status       string `json:"status"`
}
