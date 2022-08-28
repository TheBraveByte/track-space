package data

import (
	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackSpaceDBRepo interface {
	InsertInfo(email, password, id string) (int64, error)
	UpdateUserInfo(user model.User, id interface{}, t1, t2 string) error
	UpdateUserField(id, v1, v2 string) error
	VerifyLogin(id, hashedPassword, postPassword string) (bool, string)
	SendUserDetails(id string) (primitive.M, error)
	StoreWorkSpaceData(id interface{}, project model.Project) error
	ModifyProjectData(id string, project model.Project) error
	StoreDailyTaskData(task model.DailyTask, id string) error
	GetProjectData(id string) (primitive.M, error)
}
