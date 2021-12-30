package main

import (
	"GO/trevor/bookings_prj/internal/config"
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		fmt.Printf("type is %T\n", v)
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, but is %T", v))
	}
}
