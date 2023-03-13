package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"auth/data"
)

const (
	webPort = "80"
)

var (
	dbConnAttempts = 10
)

type Service struct {
	DB     *sqlx.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect to DB
	dbConn := connectToDB()
	if dbConn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up app
	service := Service{
		DB:     dbConn,
		Models: data.New(dbConn),
	}

	// create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: service.router(),
	}

	// listen and serve server
	log.Panic(srv.ListenAndServe())
}

// openDB
func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB
func connectToDB() *sqlx.DB {
	dsn := os.Getenv("DSN")
	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			dbConnAttempts--
		} else {
			log.Println("Connected to Postgres!")
			return db
		}

		if dbConnAttempts == 0 {
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(time.Second * 2)
	}
}
