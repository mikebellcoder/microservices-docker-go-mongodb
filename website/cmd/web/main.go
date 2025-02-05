package main

import "log"

type apis struct {
	users     string
	movies    string
	showtimes string
	bookings  string
}

type application struct {
	errLog  *log.Logger
	infoLog *log.Logger
	apis    apis
}
