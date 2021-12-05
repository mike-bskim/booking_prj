package main

import (
	"GO/trevor/bookings-31/pkg/config"
	"GO/trevor/bookings-31/pkg/handlers"
	"GO/trevor/bookings-31/pkg/render"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {

	// change this to true wehn in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // middleware.go > Secure:   false, 수정시 둘다 수정 필요.
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false // false: (DEV mode)read cache everytime.

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	tmp := fmt.Sprintf("Staring application on port %s", portNumber)
	fmt.Println(tmp)
	// _ = http.ListenAndServe(portNumber, nil)

	svr := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = svr.ListenAndServe()
	log.Fatal(err)
}
