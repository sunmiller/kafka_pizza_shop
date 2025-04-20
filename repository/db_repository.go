package repository

import (
	"context"
	"log"

	"github.com/sunmiller/pizza-shop-eda/order-service/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	Create(data interface{}, ctx interface{}) (interface{}, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func getSessionContext(sessionContext interface{}) mongo.SessionContext {
	cont := context.Background()
	if sessionContext == nil {
		return mongo.NewSessionContext(cont, mongo.SessionFromContext(cont))
	}
	return sessionContext.(mongo.SessionContext)
}

func (mr *MongoRepository) Create(data interface{}, ctx interface{}) (interface{}, error) {
	sc := getSessionContext(ctx)
	result, err := mr.collection.InsertOne(sc, data)
	return result, err
}

func GetMongoRepository(dbName, collectionName string) *MongoRepository {
	collection := config.GetDatabaseCollection(&dbName, collectionName)
	log.Printf("mongo collection: %v", collection.Name())
	return &MongoRepository{
		collection: collection,
	}
}
