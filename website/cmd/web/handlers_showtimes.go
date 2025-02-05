package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	moviesModel "github.com/mikebellcoder/microservices-docker-go-mongodb/movies/pkg/models"
	"github.com/mikebellcoder/microservices-docker-go-mongodb/showtimes/pkg/models"
)

type showtimeTemplateData struct {
	ShowTime  models.ShowTime
	ShowTimes []models.ShowTime
	Movies    string
}

func (app *application) showtimesList(w http.ResponseWriter, r *http.Request) {

	// get showtimes list from api
	app.infoLog.Println("Calling showtimes API...")
	resp, err := http.Get(app.apis.showtimes)
	if err != nil {
		app.infoLog.Println(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		app.infoLog.Println(err.Error())
	}

	var td showtimeTemplateData
	json.Unmarshal(bodyBytes, &td)
	app.infoLog.Println(td.ShowTimes)

	// load template files
	files := []string{
		"./ui/html/showtimes/list.page.tmpl",
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

func (app *application) showtimesView(w http.ResponseWriter, r *http.Request) {
	// get id from url
	vars := mux.Vars(r)
	showtimeID := vars["id"]

	// get showtimes from API
	app.infoLog.Println("Calling showtimes API...")
	url := fmt.Sprintf("%s/%s", app.apis.showtimes, showtimeID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var td showtimeTemplateData
	json.Unmarshal(bodyBytes, &td)
	app.infoLog.Println(td.ShowTime)

	// load movies
	var movies []string
	for _, m := range td.ShowTime.Movies {
		url := fmt.Sprintf("%s/%s", app.apis.movies, m)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err.Error())
		}

		var movie moviesModel.Movie
		json.Unmarshal(bodyBytes, &movie)
		movies = append(movies, movie.Title)
	}
	td.Movies = strings.Join(movies, ", ")

	// load template
	files := []string{
		"./ui/html/showtimes/view.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, td)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
