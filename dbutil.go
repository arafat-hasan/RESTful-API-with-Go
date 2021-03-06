package main

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatastore struct {
	db      *mongo.Database
	Client  *mongo.Client
	Context context.Context
	logger  *logrus.Logger
}

func insertOne(mongoDataStore *MongoDatastore, _collection string, doc interface{}) (*mongo.InsertOneResult, error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db

	collection := dataBase.Collection(_collection)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func insertMany(mongoDataStore *MongoDatastore, _collection string, docs []interface{}) (*mongo.InsertManyResult, error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db
	collection := dataBase.Collection(_collection)
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

func query(mongoDataStore *MongoDatastore, _collection string, query interface{}) (result *mongo.Cursor, err error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db
	collection := dataBase.Collection(_collection)
	// result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	result, err = collection.Find(ctx, query)
	return result, err
}

func UpdateOne(mongoDataStore *MongoDatastore, _collection string, filter, update interface{}) (result *mongo.UpdateResult, err error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db
	collection := dataBase.Collection(_collection)
	result, err = collection.UpdateOne(ctx, filter, update)
	return result, err
}

func UpdateMany(mongoDataStore *MongoDatastore, _collection string, filter, update interface{}) (result *mongo.UpdateResult, err error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db
	collection := dataBase.Collection(_collection)
	result, err = collection.UpdateMany(ctx, filter, update)
	return
}

func deleteOne(mongoDataStore *MongoDatastore, _collection string, query interface{}) (result *mongo.DeleteResult, err error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db
	collection := dataBase.Collection(_collection)
	result, err = collection.DeleteOne(ctx, query)
	return result, err
}

func deleteMany(mongoDataStore *MongoDatastore, _collection string, query interface{}) (result *mongo.DeleteResult, err error) {
	//client := mongoDataStore.Client
	ctx := mongoDataStore.Context
	dataBase := mongoDataStore.db
	collection := dataBase.Collection(_collection)
	result, err = collection.DeleteMany(ctx, query)
	return
}

const CONNECTED = "Successfully connected to database: "

func NewDatastore(config Configurations, logger *logrus.Logger) *MongoDatastore {

	var mongoDataStore *MongoDatastore
	db, client, ctx := connect(config, logger)
	if db != nil && client != nil {

		// log statements here as well

		mongoDataStore = new(MongoDatastore)
		mongoDataStore.db = db
		mongoDataStore.logger = logger
		mongoDataStore.Client = client
		mongoDataStore.Context = ctx
		return mongoDataStore
	}

	logger.Fatalf("Failed to connect to database: %v", config.Database.DBName)

	return nil
}

func connect(generalConfig Configurations, logger *logrus.Logger) (a *mongo.Database, b *mongo.Client, c context.Context) {
	var connectOnce sync.Once
	var db *mongo.Database
	var client *mongo.Client
	var ctx context.Context
	connectOnce.Do(func() {
		db, client, ctx = connectToMongo(generalConfig, logger)
	})

	return db, client, ctx
}

func connectToMongo(generalConfig Configurations, logger *logrus.Logger) (a *mongo.Database, b *mongo.Client, c context.Context) {

	var err error
	client, err := mongo.NewClient(options.Client().ApplyURI(generalConfig.Database.DBURI))
	if err != nil {
		logger.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client.Connect(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	var DB = client.Database(generalConfig.Database.DBName)
	logger.Info(CONNECTED, generalConfig.Database.DBName)

	return DB, client, ctx
}
