package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sunmiller/pizza-shop-eda/order-service/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBClient      *mongo.Client
	MONGO_DB_NAME = "pizza-shop-eda"
)

func init() {
	InitializeDB()
}

func InitializeDB() {
	logger.Log("initializeing database once more")
	if DBClient == nil {
		DBClient, err := initDatabase()
		if err != nil {
			log.Fatalf("failed to initialize database %v", err)
		}
		logger.Log(fmt.Sprintf("DBClient timeout: %v", DBClient.Timeout()))
	}
}

func initDatabase() (*mongo.Client, error) {
	dbURL := GetEnvProperty("DatabaseURL")
	if dbURL == "" {
		return nil, fmt.Errorf("DatabaseURL is not set int the env variable")
	}

	clientOptions := options.Client().ApplyURI(dbURL).
		SetMaxPoolSize(600).
		SetMinPoolSize(50).
		SetMaxConnIdleTime(30 * time.Second)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with mongodb: %v", err)
	}

	log.Printf("Connected to Mongo db: %s\n", dbURL)

	return client, nil
}

func GetDatabaseCollection(dbName *string, collectionName string) *mongo.Collection {
	if dbName == nil {
		dbName = &MONGO_DB_NAME
	}
	client, err := initDatabase()
	if err != nil {
		log.Fatalf("Mongodb client initialization failed: %v", err)
	}
	if client == nil {
		log.Fatalf("Mongodb client initialization failed: %v", err)
	}
	return client.Database(*dbName).Collection(collectionName)
}

func GetMongoClient() *mongo.Client {
	client, err := initDatabase()
	if err != nil || client == nil {
		log.Fatalf("Mongodb client initialization failed: %v", err)
	}
	return client
}
