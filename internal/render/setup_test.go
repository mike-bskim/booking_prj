package render

import (
	"GO/trevor/bookings_prj/internal/config"
	"GO/trevor/bookings_prj/internal/models"
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// what am I doing to put in the session
	gob.Register(models.Reservation{})

	// change this to true wehn in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction // middleware.go > Secure:   false, 수정시 둘다 수정 필요.
	testApp.Session = session

	app = &testApp

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

// type myHandler struct{}

// func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// }
