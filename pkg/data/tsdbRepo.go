package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrackSpaceDBRepo interface {
	InsertInfo(email, password string) (int64, error)
	UpdateUserInfo(info map[string]interface{}, email interface{}, t1, t2 string) error
	VerifyLogin(email string) (bool, string)
	SendUserDetails(email interface{})(primitive.M, error)
	StoreWorkSpaceData (email interface{}, projectData map[string]interface{}) error
}
