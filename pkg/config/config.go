package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

//Application config
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
}
