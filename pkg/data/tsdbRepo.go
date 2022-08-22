package data

import (
	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackSpaceDBRepo interface {
	InsertInfo(email, password string) (int64, error)
	UpdateUserInfo(user model.User, email interface{}, t1, t2 string) error
	UpdateUserField(email, v1, v2 string) error
	VerifyLogin(email string) (bool, string)
	SendUserDetails(email string) (primitive.M, error)
	StoreWorkSpaceData(email interface{}, project model.Project) error
	OrganizeWorkSpaceData(email string) (map[string]int, error)
	StoreDailyTaskData(task model.DailyTask, email string) error
	GetProjectData(id primitive.ObjectID) (primitive.M, error)
}
