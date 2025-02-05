package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mikebellcoder/microservices-docker-go-mongodb/users/pkg/models"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	// get all
	list, err := app.users.All()
	if err != nil {
		app.serverError(w, err)
	}
	// conver to JSON
	b, err := json.Marshal(list)
	if err != nil {
		app.serverError(w, err)
	}
	app.infoLog.Println("Users have been listed")
	// return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// get id
	vars := mux.Vars(r)
	id := vars["id"]

	// find by id
	user, err := app.users.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("User not found")
			return
		}
		app.serverError(w, err)
	}
	// convert to json
	b, err := json.Marshal(user)
	if err != nil {
		app.serverError(w, err)
	}
	app.infoLog.Println("Found user")

	// return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// read from request
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	// insert
	insertResult, err := app.users.Insert(u)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New user has been created, id=%s\n", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// get id
	vars := mux.Vars(r)
	id := vars["id"]

	// delete
	deleteResult, err := app.users.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("%d user(s) have been deleted\n", deleteResult.DeletedCount)
}
