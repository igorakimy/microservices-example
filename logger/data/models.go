package data

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client

	mongoDbName        = os.Getenv("MONGO_DBNAME")
	mongoQueryTimeout  = 15
	logsCollectionName = "logs"
)

type Models struct {
	LogEntry LogEntry
}

func New(mongoClient *mongo.Client) Models {
	client = mongoClient
	return Models{
		LogEntry: LogEntry{},
	}
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (le *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database(mongoDbName).Collection(logsCollectionName)
	_, err := collection.InsertOne(context.TODO(), bson.D{
		{"name", entry.Name},
		{"data", entry.Data},
		{"created_at", time.Now()},
		{"updated_at", time.Now()},
	})
	if err != nil {
		log.Printf("Error inserting into logs: %v\n", err)
		return err
	}
	return nil
}

func (le *LogEntry) GetAll() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(mongoQueryTimeout)*time.Second,
	)
	defer cancel()

	collection := client.Database(mongoDbName).Collection(logsCollectionName)

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Printf("Finding all logs error: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		if err := cursor.Decode(&item); err != nil {
			log.Printf("Error decoding log into slice: %v\n", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (le *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(mongoQueryTimeout)*time.Second,
	)
	defer cancel()

	collection := client.Database(mongoDbName).Collection(logsCollectionName)

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var logEntry *LogEntry

	res := collection.FindOne(ctx, bson.M{"_id": docId})
	if err := res.Decode(&logEntry); err != nil {
		log.Printf("Error decoding log: %v\n", err)
		return nil, err
	}

	return logEntry, nil
}

func (le *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(mongoQueryTimeout)*time.Second,
	)
	defer cancel()

	collection := client.Database(mongoDbName).Collection(logsCollectionName)

	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (le *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(mongoQueryTimeout)*time.Second,
	)
	defer cancel()

	collection := client.Database(mongoDbName).Collection(logsCollectionName)

	docID, err := primitive.ObjectIDFromHex(le.ID)
	if err != nil {
		return nil, err
	}

	res, err := collection.UpdateByID(ctx, docID, bson.D{{"$set", bson.D{
		{"name", le.Name},
		{"data", le.Data},
		{"updated_at", time.Now()},
	}}})
	if err != nil {
		return nil, err
	}

	return res, nil
}
