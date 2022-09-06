package data

import (
	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackSpaceDBRepo interface {
	InsertUserInfo(id, email, password string) (int64, error)
	UpdateUserInfo(user model.User, id string, t1, t2 string) error
	UpdateUserField(id, t1, t2 string) error
	VerifyLogin(id, hashedPassword, postPassword string) (bool, string)
	SendUserDetails(id string) (primitive.M, error)
	StoreProjectData(id string, project model.Project) error
	GetProjectData(project_id string) (primitive.M, error)
	ModifyProjectData(id string, project model.Project) error

	StoreDailyTaskData(task model.DailyTask, id string) error
	
	UpdateUserStat(data model.Data, id string) error
	GetUserStatByID(id string) (primitive.M, error)
	
	DeleteUserProject(project_id string) error
}
