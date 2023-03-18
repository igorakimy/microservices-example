package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PORT = "80"
)

type Service struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	service := Service{
		Rabbit: rabbitConn,
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: service.router(),
	}

	fmt.Printf("Starting broker on port: %s\n", PORT)

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Serving error: %v\n", err)
		panic(err)
	}
}

// connect tries connecting to RabbitMQ  and returns connection.
func connect() (*amqp.Connection, error) {
	var counts int64
	var connAttempts = 5
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(os.Getenv("RMQ_URL"))
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > int64(connAttempts) {
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
