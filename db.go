package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB abstracts database clients
type DB struct {
	MongoClient *mongo.Database
}

// NewDB is a constructor for initializing the database connections
func NewDB() *DB {
	return &DB{
		MongoClient: initMongo(),
	}
}

// called from main, connect to mongo
func initMongo() *mongo.Database {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	if mongoHost == "" || mongoPort == "" {
		log.Fatal("error MONGO_HOST or MONGO_PORT does not exist.")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)))
	if err != nil {
		log.Fatal("error creating a mongo client: ", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal("error connecting to mongodb: ", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("error pinging the mongo server: ", err)
	}

	return client.Database("backend-homework")
}

// PopulateDatabase sets up sample users and sample "likes" between the users
func PopulateDatabase(db *DB) {
	usersColl := db.MongoClient.Collection("users")
	ctx := context.Background()

	// if no users in db, add defaults
	c, err := usersColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("error counting users from mongo: ", err)
	}

	if c != 0 {
		return
	}

	users := []*User{
		{
			Age:         30,
			Bio:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			CreatedDate: time.Now(),
			ID:          primitive.NewObjectID().Hex(),
			JobTitle:    "Software Engineer",
			Name:        "Jennifer",
		},
		{
			Age:         43,
			Bio:         "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat",
			CreatedDate: time.Now(),
			ID:          primitive.NewObjectID().Hex(),
			JobTitle:    "Musician",
			Name:        "Bob",
		},
		{
			Age:         22,
			Bio:         "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			CreatedDate: time.Now(),
			ID:          primitive.NewObjectID().Hex(),
			JobTitle:    "Professor",
			Name:        "Susan",
		},
		{
			Age:         27,
			Bio:         "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			CreatedDate: time.Now(),
			ID:          primitive.NewObjectID().Hex(),
			JobTitle:    "Professional Dancer",
			Name:        "Michael",
		},
		{
			Age:         35,
			Bio:         "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			CreatedDate: time.Now(),
			ID:          primitive.NewObjectID().Hex(),
			JobTitle:    "Accountant",
			Name:        "Alexis",
		},
		{
			Age:         38,
			Bio:         "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			CreatedDate: time.Now(),
			ID:          primitive.NewObjectID().Hex(),
			JobTitle:    "Security Officer",
			Name:        "Andrew",
		},
	}

	ui := make([]interface{}, 0, len(users))

	for _, v := range users {
		ui = append(ui, v)
	}

	if _, err := usersColl.InsertMany(ctx, ui); err != nil {
		log.Fatal("error inserting users: ", err)
	}

	// Bob like Susan
	// Bob likes Michael
	// Alexis likes Michael
	// Alexis likes Andrew
	// {Everyone} likes Jennifer
	likes := []*Rating{
		{
			CreatedDate: time.Now(),
			FromUserID:  users[1].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[2].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[3].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[2].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[4].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[3].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[4].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[5].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[1].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[0].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[2].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[0].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[3].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[0].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[4].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[0].ID,
			Type:        LIKE,
		},
		{
			CreatedDate: time.Now(),
			FromUserID:  users[5].ID,
			ID:          primitive.NewObjectID().Hex(),
			ToUserID:    users[0].ID,
			Type:        LIKE,
		},
	}

	li := make([]interface{}, 0, len(likes))

	for _, v := range likes {
		li = append(li, v)
	}

	ratingsColl := db.MongoClient.Collection("ratings")
	if _, err := ratingsColl.InsertMany(ctx, li); err != nil {
		log.Fatal("error inserting likes: ", err)
	}
}
