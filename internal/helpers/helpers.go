package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bartoszjasak/bookings/internal/config"
)

var appConfig *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	appConfig = a
}

func ClientError(w http.ResponseWriter, status int) {
	appConfig.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	appConfig.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := appConfig.Session.Exists(r.Context(), "user_id")
	return exists
}
