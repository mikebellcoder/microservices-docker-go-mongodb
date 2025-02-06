package main

import "github.com/gorilla/mux"

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", app.home)
	// users
	r.HandleFunc("/users/list", app.usersList)
	r.HandleFunc("/users/view/{id}", app.usersView)
	// movies
	r.HandleFunc("/movies/list", app.moviesList)
	r.HandleFunc("/movies/view/{id}", app.moviesView)
	// showtimes
	r.HandleFunc("/showtimes/list", app.moviesList)
	r.HandleFunc("/showtimes/view/{id}", app.moviesView)
	// bookings
	r.HandleFunc("/bookings/list", app.bookingsList)
	r.HandleFunc("/bookings/view/{id}", app.bookingsView)

	// serve static resources
	r.PathPrefix("/static/").Handler(app.static("./ui/static/"))

	return r
}
