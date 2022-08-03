package tsRepoStore

import (
	"github.com/yusuf/track-space/pkg/config"
	"github.com/yusuf/track-space/pkg/data"
	"go.mongodb.org/mongo-driver/mongo"
)

type TsMongoDBRepo struct {
	AppConfig *config.AppConfig
	TsMongoDB *mongo.Client
}

func NewTsMongoDBRepo(app *config.AppConfig, tsm *mongo.Client) data.TrackSpaceDBRepo {
	return &TsMongoDBRepo{
		AppConfig: app,
		TsMongoDB: tsm,
	}

}
