package data

import (
	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TrackSpaceDBRepo : interface for all the database queries
type TrackSpaceDBRepo interface {
	// Queries for user to interact with the database
	InsertUserInfo(id, email, password string) (int64, error)
	UpdateUserInfo(user model.User, id string, t1, t2 string) error
	UpdateUserField(id, t1, t2 string) error
	VerifyLogin(id, hashedPassword, postPassword string) (bool, string)
	SendUserDetails(id string) (primitive.M, error)

	// Queries for User Project
	StoreProjectData(id string, project model.Project) error
	GetProjectData(project_id string) (primitive.M, error)
	ModifyProjectData(user_id string, id string, project model.Project) error

	// Queries for User Todo Task
	StoreTodoData(todo model.Todo, id string) error
	GetTodoData(todo_id string) (primitive.M, error)
	ModifyTodoData(id string, todo model.Todo) error

	// Queries for User Statistics
	UpdateUserStat(data model.Data, id string) error
	GetUserStatByID(id string) (primitive.M, error)

	// Queries for User to Delete Project and Todo Task
	DeleteUserProject(project_id string) error
	DeleteUserTodo(todo_id string) error

	// Queries for Admin
	GetAllUserData() ([]primitive.M, error)
	GetAdminInfo() ([]primitive.M, error)
	AdminDeleteUserData(id string) error
}
