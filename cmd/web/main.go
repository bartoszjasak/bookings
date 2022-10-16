package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bartoszjasak/bookings/pkg/config"
	"github.com/bartoszjasak/bookings/pkg/handlers"
	"github.com/bartoszjasak/bookings/pkg/render"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {

	appConfig.InProduction = false
	appConfig.UseCache = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction
	appConfig.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	appConfig.TemplateCatche = templateCache

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)

	fmt.Println(fmt.Sprintf("Starting an application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}