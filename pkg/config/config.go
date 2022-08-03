package config

import (
	"github.com/gin-contrib/sessions"
	"html/template"
	"log"
)

type AppConfig struct {
	InfoLogger      *log.Logger
	ErrorLogger     *log.Logger
	AppInProduction bool
	UseTempCache    bool
	TsData          sessions.Session
	TemplateCache   map[string]*template.Template
}
