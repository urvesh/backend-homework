package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Like is a relationship struct where it contains which user likes which other user
type Like struct {
	ID          string    `json:"_id" bson:"_id"`
	UserID      string    `json:"userId" bson:"userId"`
	LikeUserID  string    `json:"likeUserId" bson:"likeUserId"`
	CreatedDate time.Time `json:"createdDate" bson:"createdDate"`
}

// look for all the users that liked this userId
func FindLikesByUserID(db *DB, userId string) ([]*Like, error) {
	coll := db.MongoClient.Collection("likes")
	ctx := context.Background()

	likes := make([]*Like, 0)

	filter := bson.M{
		"likeUserId": userId,
	}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		log.Println("error finding likes from mongo: ", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var l Like
		if err := cur.Decode(&l); err != nil {
			log.Println("error decoding into user struct", err)
			return nil, err
		}

		likes = append(likes, &l)
	}

	return likes, nil
}

// Save inserts a new like entry to the database
func (l *Like) Save(db *DB) error {
	coll := db.MongoClient.Collection("likes")
	ctx := context.Background()

	l.ID = primitive.NewObjectID().Hex()
	l.CreatedDate = time.Now()

	_, err := coll.InsertOne(ctx, l)
	if err != nil {
		log.Printf("error creating a new like: %s\n", err)
		return err
	}

	return nil
}
