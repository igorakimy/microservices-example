package main

import (
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	Mailer *Mail
}

const (
	webPort = "80"
)

func main() {
	service := Service{
		Mailer: NewMail(),
	}

	log.Printf("Starting mail service on port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: service.router(),
	}

	log.Panic(srv.ListenAndServe())
}
