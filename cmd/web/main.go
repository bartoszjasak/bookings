package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bartoszjasak/bookings/internal/config"
	"github.com/bartoszjasak/bookings/internal/handlers"
	"github.com/bartoszjasak/bookings/internal/helpers"
	"github.com/bartoszjasak/bookings/internal/models"
	"github.com/bartoszjasak/bookings/internal/render"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Starting an application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	gob.Register(models.Reservation{})

	appConfig.InProduction = false
	appConfig.UseCache = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction
	appConfig.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	appConfig.TemplateCatche = templateCache

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig)
	helpers.NewHelpers(&appConfig)
	return nil
}
