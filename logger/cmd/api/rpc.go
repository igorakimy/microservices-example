package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

// RPCServer is the type for our RPC Server. Methods that take this as
// a receiver are available over RPC, as long as they are exported.
type RPCServer struct{}

// RPCPayload is the type for data we receive from RPC.
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo writes our payload to mongo.
func (rs *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database(data.GetMongoDBName()).
		Collection(data.GetCollectionName())
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("Error writing to mongo: %v\n", err)
		return err
	}

	*resp = "Processed payload via RPC: " + payload.Name
	return nil
}
