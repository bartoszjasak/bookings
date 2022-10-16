package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bartoszjasak/bookings/pkg/config"
	"github.com/bartoszjasak/bookings/pkg/models"
)

var appConfig *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

func RenderTemplate(w http.ResponseWriter, tmplName string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template
	if appConfig.UseCache {
		templateCache = appConfig.TemplateCatche
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	template, ok := templateCache[tmplName]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	_ = template.Execute(buf, templateData)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	var templateCache = map[string]*template.Template{}

	pages, err := filepath.Glob("templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		template, err := template.ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		layouts, err := filepath.Glob("templates/*.layout.tmpl")
		if err != nil {
			return templateCache, err
		}

		if len(layouts) > 0 {
			template, err = template.ParseGlob("templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = template
	}
	return templateCache, nil
}
