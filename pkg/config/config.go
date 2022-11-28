package config

import (
	"github.com/go-playground/validator/v10"
	"log"

	"github.com/yusuf/track-space/pkg/model"
)

type AppConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	MailChan    chan model.Email
	Validator   *validator.Validate
}
