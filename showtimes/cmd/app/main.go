package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikebellcoder/microservices-docker-go-mongodb/showtimes/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errLog    *log.Logger
	infoLog   *log.Logger
	showtimes *mongodb.ShowTimeModel
}

func main() {

	serverAddr := flag.String("serverAddr", "", "HTTP server network address")
	serverPort := flag.Int("serverPort", 4000, "HTTP server network port")
	mongoURI := flag.String("mongoURI", "mongodb://localhost:27017", "Database hostname url")
	mongodDatabase := flag.String("mongoDatabase", "showtimes", "Database name")
	enableCredentials := flag.Bool("enableCredentials", false, "Enable the use of credentials for mongo connection")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	co := options.Client().ApplyURI(*mongoURI)
	if *enableCredentials {
		co.Auth = &options.Credential{
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, co)
	if err != nil {
		errLog.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	infoLog.Printf("Database connection established")

	app := &application{
		infoLog: infoLog,
		errLog:  errLog,
		showtimes: &mongodb.ShowTimeModel{
			C: client.Database(*mongodDatabase).Collection("showtimes"),
		},
	}

	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	srv := &http.Server{
		Addr:         serverURI,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", serverURI)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
