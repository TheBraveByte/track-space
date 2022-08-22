package config

import (
	"html/template"
	"log"

	"github.com/gin-contrib/sessions"
)

type AppConfig struct {
	InfoLogger      *log.Logger
	ErrorLogger     *log.Logger
	AppInProduction bool
	UseTempCache    bool
	TsData          sessions.Session
	TemplateCache   map[string]*template.Template
}
