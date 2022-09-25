package tsRepoStore

import (
	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/data"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
TsMongoDBRepo Database setup for controller package
*/
type TsMongoDBRepo struct {
	AppConfig *config.AppConfig
	TsMongoDB *mongo.Client
}

/*
NewTsMongoDBRepo : function to keep track or update on request for executing call database
query method
*/
func NewTsMongoDBRepo(app *config.AppConfig, tsm *mongo.Client) data.TrackSpaceDBRepo {
	return &TsMongoDBRepo{
		AppConfig: app,
		TsMongoDB: tsm,
	}
}
