package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/bartoszjasak/bookings/internal/config"
	"github.com/bartoszjasak/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
}

var appConfig *config.AppConfig
var pathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	appConfig = a
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.Flash = appConfig.Session.PopString(r.Context(), "flash")
	templateData.Error = appConfig.Session.PopString(r.Context(), "error")
	templateData.Warning = appConfig.Session.PopString(r.Context(), "warning")
	templateData.CSRFToken = nosurf.Token(r)
	if appConfig.Session.Exists(r.Context(), "user_id") {
		templateData.IsAuthenticated = 1
	}
	return templateData
}

func Template(w http.ResponseWriter, r *http.Request, tmplName string, templateData *models.TemplateData) error {
	var templateCache map[string]*template.Template
	if appConfig.UseCache {
		templateCache = appConfig.TemplateCatche
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	template, ok := templateCache[tmplName]
	if !ok {
		return errors.New("cannot get template from cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, r)

	_ = template.Execute(buf, templateData)

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	var templateCache = map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		template, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return templateCache, err
		}

		if len(layouts) > 0 {
			template, err = template.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = template
	}
	return templateCache, nil
}
