package databases

import (
	"context"

	"api.workzen.odoo/constants"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient is the MongoDB client
var MongoDBClient *mongo.Client

// MongoDBDatabase is the MongoDB database
var MongoDBDatabase *mongo.Database

// ConnectMongoDB connects to MongoDB
func ConnectMongoDB() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(constants.DatabaseMongodbURI))
	if err != nil {
		return err
	}

	MongoDBClient = client
	MongoDBDatabase = client.Database(constants.DatabaseMongodbDBName)

	if err := MongoDBClient.Ping(context.Background(), nil); err != nil {
		return err
	}

	return nil
}

// DisconnectMongoDB disconnects from MongoDB
func DisconnectMongoDB() error {
	if err := MongoDBClient.Disconnect(context.Background()); err != nil {
		return err
	}

	return nil
}

func GetMongoDBDatabase() *mongo.Database {
	return MongoDBDatabase
}

func GetMongoDBCollection(collectionName string) *mongo.Collection {
	return MongoDBDatabase.Collection(collectionName)
}
