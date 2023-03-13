package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	PORT = "80"
)

type Service struct{}

func main() {
	var service Service

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: service.router(),
	}

	fmt.Printf("Starting broker on port: %s\n", PORT)

	log.Panic(srv.ListenAndServe())
}
