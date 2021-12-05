package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

// NoSurf andds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	log.Println("call NoSurf")
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteDefaultMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every requests
func SessionLoad(next http.Handler) http.Handler {
	log.Println("call SessionLoad")
	return session.LoadAndSave(next)
}
