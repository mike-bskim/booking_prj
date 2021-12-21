package main

import (
	// "GO/trevor/bookings_prj/pkg/config"
	// "GO/trevor/bookings_prj/pkg/handlers"
	// "GO/trevor/bookings_prj/pkg/render"
	"GO/trevor/bookings_prj/internal/config"
	"GO/trevor/bookings_prj/internal/handlers"
	"GO/trevor/bookings_prj/internal/models"
	"GO/trevor/bookings_prj/internal/render"
	"encoding/gob"
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

	err := run()
	if err != nil {
		log.Fatal(err)
	}

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

func run() error {

	// what am I doing to put in the session
	gob.Register(models.Reservation{})

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
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false // false: (DEV mode)read cache everytime.

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	return nil
}
