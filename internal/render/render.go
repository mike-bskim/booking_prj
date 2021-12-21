package render

import (
	"GO/trevor/bookings_prj/internal/config"
	"GO/trevor/bookings_prj/internal/models"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
	// 아래 부분이 html 에서 불러오는 이름("{{.CSRFToken}}")과 같아야 함.
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// DEV mode 에서는 UseCache==false 이므로. read cache everytime.
		tc, _ = CreateTemplateCache()
	}

	// map에 원하는 페이지가 있는지 확인
	t, ok := tc[tmpl]
	if !ok {
		// return errors.New("Could not get template from template cache")
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer) // buf 생성
	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td) // 해당 페이지를 buf 에 저장
	// if err != nil {
	// 	return err
	// }
	_, err := buf.WriteTo(w) // client 에게 전송
	if err != nil {
		// log.Println(err)
		fmt.Println("error writing template to browser", err)
		// return err
	}

}

// CreateTemplateCache creates a template cache as a map, page + base
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	// 폴더이름+파일명 저장
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		// 폴더정보를 제거하고 파일 이름만 저장
		name := filepath.Base(page)

		// 페이지 정보 로딩
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(layouts) > 0 {
			// 페이지 정보에 base 정보 결함.
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}

		}
		myCache[name] = ts

	}

	return myCache, nil
}
