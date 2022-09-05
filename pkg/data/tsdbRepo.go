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
	StoreWorkSpaceData(id string, project model.Project) error
	ModifyProjectData(id string, project model.Project) error
	StoreDailyTaskData(task model.DailyTask, id string) error
	GetProjectData(project_id string) (primitive.M, error)
	UpdateProjectStat(data model.Data, id string) error
	GetProjectStatByID(id string) (primitive.M, error)
}
