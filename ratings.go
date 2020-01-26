package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LIKE   = "LIKE"
	BLOCK  = "BLOCK"
	REPORT = "REPORT"
)

// Rating is a struct where it contains which user likes/blocks/reports other users
type Rating struct {
	CreatedDate time.Time `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
	FromUserID  string    `json:"fromUserId,omitempty" bson:"fromUserId,omitempty"`
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Reason      string    `json:"reason,omitempty" bson:"reason,omitempty"`
	ToUserID    string    `json:"toUserId,omitempty" bson:"toUserId,omitempty"`
	Type        string    `json:"type,omitempty" bson:"type,omitempty"`
}

type RatingParams struct {
	Filter Rating
	// Projection
}

// look for all the users that has took an action against this userId
func FindRatings(db *DB, params RatingParams) ([]*Rating, error) {
	coll := db.MongoClient.Collection("ratings")
	ctx := context.Background()

	ratings := make([]*Rating, 0)

	cur, err := coll.Find(ctx, params.Filter)
	if err != nil {
		log.Println("error finding ratings from mongo: ", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var r Rating
		if err := cur.Decode(&r); err != nil {
			log.Println("error decoding into user struct", err)
			return nil, err
		}

		ratings = append(ratings, &r)
	}

	return ratings, nil
}

// Save inserts a new like entry to the database
func (r *Rating) Save(db *DB) error {
	coll := db.MongoClient.Collection("ratings")
	ctx := context.Background()

	r.ID = primitive.NewObjectID().Hex()
	r.CreatedDate = time.Now()

	// TODO: make sure that like doesn't already exist...

	_, err := coll.InsertOne(ctx, r)
	if err != nil {
		log.Printf("error creating a new like: %s\n", err)
		return err
	}

	return nil
}
