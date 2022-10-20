package main

import (
	"fmt"
	"testing"

	"github.com/bartoszjasak/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var appConfig config.AppConfig

	mux := routes(&appConfig)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("Type is not *chi.Mux, but: %T", v))
	}
}
