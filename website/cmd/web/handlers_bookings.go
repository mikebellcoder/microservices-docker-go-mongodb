package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/mikebellcoder/microservices-docker-go-mongodb/bookings/pkg/models"
	modelsShowTime "github.com/mikebellcoder/microservices-docker-go-mongodb/showtimes/pkg/models"
	modelsUser "github.com/mikebellcoder/microservices-docker-go-mongodb/users/pkg/models"
)

type bookingTemplateData struct {
	Booking      models.Booking
	Bookings     []models.Booking
	BookingData  bookingData
	BookingsData []bookingData
}

type bookingData struct {
	ID           string
	UserFullName string
	ShowTimeDate string
}

func (app *application) loadBookingData(btd *bookingTemplateData, isList bool) {
	// clean booking data
	btd.BookingsData = []bookingData{}
	btd.BookingData = bookingData{}

	// load booking data
	if isList {
		for _, b := range btd.Bookings {
			// load user data
			userURL := fmt.Sprintf("%s/%s", app.apis.users, b.UserID)
			var user modelsUser.User
			err := app.getAPIContent(userURL, &user)
			if err != nil {
				app.errLog.Println(err.Error())
			}

			// load showtime data
			showtimeURL := fmt.Sprintf("%s/%s", app.apis.showtimes, b.ShowtimeID)
			var showtime modelsShowTime.ShowTime
			err = app.getAPIContent(showtimeURL, &showtime)
			if err != nil {
				app.errLog.Println(err.Error())
			}

			bookingData := bookingData{
				ID:           b.ID.Hex(),
				UserFullName: fmt.Sprintf("%s %s", user.Name, user.LastName),
				ShowTimeDate: showtime.Date,
			}
			btd.BookingsData = append(btd.BookingsData, bookingData)
			app.infoLog.Println(b.UserID)
		}
	} else {
		b := btd.Booking

		// load user data
		userURL := fmt.Sprintf("%s/%s", app.apis.users, b.UserID)
		var user modelsUser.User
		err := app.getAPIContent(userURL, &user)
		if err != nil {
			app.errLog.Println(err.Error())
		}

		// laod showtime data
		showtimeURL := fmt.Sprintf("%s/%s", app.apis.showtimes, b.ShowtimeID)
		var showtime modelsShowTime.ShowTime

		err = app.getAPIContent(showtimeURL, &showtime)
		if err != nil {
			app.errLog.Println(err.Error())
		}

		btd.BookingData = bookingData{
			ID:           b.ID.Hex(),
			UserFullName: fmt.Sprintf("%s %s", user.Name, user.LastName),
			ShowTimeDate: showtime.Date,
		}
	}
}

func (app *application) bookingsList(w http.ResponseWriter, r *http.Request) {
	// get bookings list from api
	var td bookingTemplateData
	app.infoLog.Println("Calling bookings API...")

	err := app.getAPIContent(app.apis.bookings, &td.Bookings)
	if err != nil {
		app.errLog.Println(err.Error())
	}
	app.infoLog.Println(td.Bookings)
	app.infoLog.Println(td)

	app.loadBookingData(&td, true)

	// load template files
	files := []string{
		"./ui/html/bookings/list.page.tmpl",
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

func (app *application) bookingsView(w http.ResponseWriter, r *http.Request) {
	// get id from incoming url
	vars := mux.Vars(r)
	bookingID := vars["id"]

	// get booking list from API
	var td bookingTemplateData
	app.infoLog.Println("Calling bookings API...")
	url := fmt.Sprintf("%s/%s", app.apis.bookings, &bookingID)

	err := app.getAPIContent(url, &td.Booking)
	if err != nil {
		app.errLog.Println(err.Error())
	}
	app.infoLog.Println(td.Booking)
	app.infoLog.Println(url)

	app.loadBookingData(&td, false)

	// load template files
	files := []string{
		"./ui/html/bookings/view.page.tmpl",
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
		return
	}
}
