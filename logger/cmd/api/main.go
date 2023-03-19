package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
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
	rpcHost  = os.Getenv("RPC_HOST")
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

	// Register the RPC server
	err = rpc.Register(new(RPCServer))
	go service.rpcListen()

	// Start gRPC server
	go service.gRPCListen()

	// start web server
	log.Printf("Starting service on port: %v\n", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: service.router(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Serving error: %v\n", err)
		panic(err)
	}
}

// rpcListen registers and serve RPC connection.
func (s *Service) rpcListen() error {
	log.Printf("Starting RPC server on %s:%s", rpcHost, rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", rpcHost, rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
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
