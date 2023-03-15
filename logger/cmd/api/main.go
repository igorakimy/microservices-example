package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"logger/data"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoConnTimeout = 15
)

var (
	webPort  = os.Getenv("LOGGER_SERVICE_PORT")
	rpcPort  = os.Getenv("RPC_PORT")
	mongoURI = os.Getenv("MONGO_URI")
	grpcPort = os.Getenv("GRPC_PORT")

	client *mongo.Client
)

type Service struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongoDB()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(mongoConnTimeout),
	)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Mongo client disconnect error: %v\n", err)
			panic(err)
		}
	}()

	service := Service{
		Models: data.New(client),
	}

	// start web server
	log.Printf("Starting service on port: %v\n", webPort)
	// service.serve()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: service.router(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Serving error: %v\n", err)
		panic(err)
	}
}

// serve handle incoming connections.
func (s *Service) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: s.router(),
	}

	log.Panic(srv.ListenAndServe())
}

func connectToMongoDB() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})

	// connect to mongo and get client
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("Error connection to MongoDB: %v\n", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
