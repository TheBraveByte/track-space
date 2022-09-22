package data

import (
	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TrackSpaceDBRepo : interface for all the database queries
type TrackSpaceDBRepo interface {
	// Queries for user to interact with the database

	InsertUserInfo(email, password string) (int64, string, error)
	UpdateUserInfo(user model.User, id string, t1, t2 string) error
	UpdateUserField(id, t1, t2 string) error
	VerifyLogin(id, hashedPassword, postPassword string) (bool, string)
	ResetUserPassword(email, newPassword string) error
	SendUserDetails(id string) (primitive.M, error)

	// Queries for User Project

	StoreProjectData(id string, project model.Project) error
	GetProjectData(projectId string) (primitive.M, error)
	ModifyProjectData(userId string, id string, project model.Project) error

	// Queries for User Todo Task

	StoreTodoData(todo model.Todo, id string) error
	GetTodoData(todoId string) (primitive.M, error)
	ModifyTodoData(id string, todo model.Todo) error

	// Queries for User Statistics

	UpdateUserStat(data model.Data, id string) error
	GetUserStatByID(id string) (primitive.M, error)

	// Queries for User to Delete Project and Todo Task

	DeleteUserProject(projectId string) error
	DeleteUserTodo(todoId string) error

	// Queries for Admin

	GetAllUserData() ([]primitive.M, error)
	GetAdminInfo() ([]primitive.M, error)
	AdminDeleteUserData(id string) error
}
