package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/mikebellcoder/microservices-docker-go-mongodb/users/pkg/models"
)

type userTemplateData struct {
	User  models.User
	Users []models.User
}

func (app *application) usersList(w http.ResponseWriter, r *http.Request) {
	// get users list from api
	var utd userTemplateData
	err := app.getAPIContent(app.apis.users, &utd.Users)
	if err != nil {
		app.errLog.Println(err.Error())
	}
	app.infoLog.Println(utd.Users)

	// load template files
	files := []string{
		"./ui/html/users/list.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, utd)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) usersView(w http.ResponseWriter, r *http.Request) {
	// get id from url
	vars := mux.Vars(r)
	userID := vars["id"]

	// get users list from api
	app.infoLog.Println("Calling users API...")
	url := fmt.Sprintf("%s/%s", app.apis.users, userID)

	var utd userTemplateData
	app.getAPIContent(url, &utd.User)
	app.infoLog.Println(utd.User)

	// load template files
	files := []string{
		"./ui/html/users/view.page.tmpl",
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
	err = ts.Execute(w, utd.User)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
