package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	Filter *Rating
	// Projection
}

// look for all the users that has took an action against this userId
func FindRatings(db *DB, params *RatingParams) ([]*Rating, error) {
	coll := db.MongoClient.Collection("ratings")
	ctx := context.Background()

	ratings := make([]*Rating, 0)

	cur, err := coll.Find(ctx, params.Filter)
	if err != nil {
		return nil, NewErrorf("error finding ratings from mongo: %s", err)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var r Rating
		if err := cur.Decode(&r); err != nil {
			return nil, NewErrorf("error decoding into user struct: %s", err)
		}

		ratings = append(ratings, &r)
	}

	return ratings, nil
}

// FindRatingExists checks if a given document exists within the db
func FindRatingExists(db *DB, params *RatingParams) (bool, error) {
	coll := db.MongoClient.Collection("ratings")

	doc := coll.FindOne(context.Background(), params.Filter)
	if doc.Err() != nil {
		if doc.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		log.Printf("error looking up rating: %s \n", doc.Err())
		return false, doc.Err()
	}

	return true, nil
}

// Save inserts a new like entry to the database
func (r *Rating) Save(db *DB) error {
	coll := db.MongoClient.Collection("ratings")
	ctx := context.Background()

	// create unique ID and set createdDate
	r.ID = primitive.NewObjectID().Hex()
	r.CreatedDate = time.Now()

	filter := &RatingParams{
		Filter: &Rating{
			FromUserID: r.FromUserID,
			ToUserID:   r.ToUserID,
			Type:       r.Type,
		},
	}

	// check if the rating exists before inserting a new one
	exists, err := FindRatingExists(db, filter)
	if err != nil {
		return err
	}

	// don't save new entry if it already exists.
	if exists {
		return nil
	}

	// otherwise insert a new one
	if _, err := coll.InsertOne(ctx, r); err != nil {
		log.Printf("error creating a new like: %s\n", err)
		return err
	}

	// if it was a block entry, from user A to B, remove user A's LIKE to user B, and vice versa
	if r.Type == BLOCK {
		filter.Filter.Type = LIKE
		if _, err := coll.DeleteOne(ctx, filter); err != nil {
			return NewErrorf("error deleting like entry: %s", err)
		}

		// swap fromUserId with toUserId
		filter.Filter.FromUserID, filter.Filter.ToUserID = filter.Filter.ToUserID, filter.Filter.FromUserID
		if _, err := coll.DeleteOne(ctx, filter); err != nil {
			return NewErrorf("error deleting like entry: %s", err)
		}
	}

	return nil
}
