package handlers

import (
	"fmt"
	"net/http"

	"github.com/bartoszjasak/bookings/pkg/config"
	"github.com/bartoszjasak/bookings/pkg/models"
	"github.com/bartoszjasak/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	AppConfig *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	fmt.Println(remoteIp)
	m.AppConfig.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello, again."

	remoteIP := m.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) DoNothing(w http.ResponseWriter, r *http.Request) {
}
