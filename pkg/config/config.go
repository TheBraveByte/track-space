package config

import (
	"log"

	"github.com/yusuf/track-space/pkg/model"
)

type AppConfig struct {
	InfoLogger      *log.Logger
	ErrorLogger     *log.Logger
	AppInProduction bool
	UseTempCache    bool
	MailChan        chan model.Email
}
