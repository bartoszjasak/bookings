package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bartoszjasak/bookings/internal/config"
	"github.com/bartoszjasak/bookings/internal/models"
)

var session *scs.SessionManager
var testAppConfig config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testAppConfig.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testAppConfig.ErrorLog = errorLog

	testAppConfig.UseCache = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testAppConfig.Session = session

	appConfig = &testAppConfig

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {
}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
