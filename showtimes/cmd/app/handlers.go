package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mikebellcoder/microservices-docker-go-mongodb/showtimes/pkg/models"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {
	showtimes, err := app.showtimes.All()
	if err != nil {
		app.serverError(w, err)
	}

	b, err := json.Marshal(showtimes)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Showtimes have been listed")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := app.showtimes.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Showtime not found")
			return
		}
		app.serverError(w, err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Found a showtime")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findByDate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	date := vars["date"]

	m, err := app.showtimes.FindByDate(date)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Shotime not found")
			return
		}
		app.serverError(w, err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("Date found a showtime")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	var m models.ShowTime

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	m.CreatedAt = time.Now()
	insertResult, err := app.showtimes.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("New showtime has been created, id=%s\n", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	deleteResult, err := app.showtimes.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("Deleted %d showtime(s)\n", deleteResult.DeletedCount)
}
