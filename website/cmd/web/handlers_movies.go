package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikebellcoder/microservices-docker-go-mongodb/movies/pkg/models"
)

type movieTemplateData struct {
	Movie  models.Movie
	Movies []models.Movie
}

func (app *application) moviesList(w http.ResponseWriter, r *http.Request) {
	// get movies list
	var mtd movieTemplateData
	app.infoLog.Println("Calling movies API...")
	app.getAPIContent(app.apis.movies, &mtd.Movies)
	app.infoLog.Println(mtd.Movies)

	// load template file
	files := []string{
		"./ui/html/movies/list.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// respond
	err = ts.Execute(w, mtd)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) moviesView(w http.ResponseWriter, r *http.Request) {
	// get id from url
	vars := mux.Vars(r)
	movieID := vars["id"]

	// get movies list
	app.infoLog.Println("Calling movies API...")
	url := fmt.Sprintf("%s/%s", app.apis.movies, movieID)

	resp, err := http.Get(url)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var td movieTemplateData
	json.Unmarshal(bodyBytes, &td)
	app.infoLog.Println(td.Movie)

	// load template
	files := []string{
		"./ui/html/movies/view.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// respond
	err = ts.Execute(w, td)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
