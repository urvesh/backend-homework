package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Age         int       `json:"age,omitempty" bson:"age,omitempty"`
	Bio         string    `json:"bio,omitempty" bson:"bio,omitempty"`
	CreatedDate time.Time `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
	ID          string    `json:"_id,omitempty" bson:"_id,omitempty"`
	JobTitle    string    `json:"jobTitle,omitempty" bson:"jobTitle,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty"`
}

// FindAllUsers returns all the existing users from the db
func FindAllUsers(db *DB) ([]*User, error) {
	coll := db.MongoClient.Collection("users")
	ctx := context.Background()

	users := make([]*User, 0)

	filter := bson.M{}
	cur, err := coll.Find(ctx, filter)

	if err != nil {
		log.Println("error finding users from mongo:", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var u User
		if err := cur.Decode(&u); err != nil {
			log.Println("error decoding into user struct", err)
			return nil, err
		}

		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		log.Println("mongo error", err)
		return nil, err
	}

	return users, nil
}

// helper function to do simple id lookups
func findByID(db *DB, collection, fieldName, id string) (*mongo.SingleResult, error) {
	coll := db.MongoClient.Collection(collection)
	ctx := context.Background()

	filter := bson.M{
		fieldName: id,
	}

	doc := coll.FindOne(ctx, filter)
	if doc.Err() != nil {
		log.Printf("error looking up %s %s: %s \n", collection, id, doc.Err())

		if doc.Err() == mongo.ErrNoDocuments {
			// document not found, not an actual error
			return nil, nil
		}
		return nil, doc.Err()
	}

	return doc, nil
}

// FindUserById lookup user by id
func FindUserByID(db *DB, id string) (*User, error) {
	doc, err := findByID(db, "users", "_id", id)
	if err != nil || doc == nil {
		// if it was a 404, err would also be nil
		return nil, err
	}

	var u User
	if err := doc.Decode(&u); err != nil {
		log.Printf("error decoding user %s into struct: %s \n", id, err)
		return nil, err
	}

	return &u, nil
}

// FindIncomingLikes finds all the users who have liked the given userId
func FindIncomingLikes(db *DB, userId string) ([]*User, error) {
	p := RatingParams{
		Filter: Rating{
			ToUserID: userId,
			Type:     LIKE,
		},
	}
	likes, err := FindRatings(db, p)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(likes))

	// return empty array if no likes.
	if len(likes) == 0 {
		return users, nil
	}

	// go through found likes and populate it with user data
	for _, v := range likes {
		u, err := FindUserByID(db, v.FromUserID)
		if err != nil {
			return nil, err
		}

		if u != nil {
			users = append(users, u)
		}
	}

	return users, nil
}

// Edit overrides user data with the incoming values
func (u *User) Edit(db *DB) (*User, error) {
	coll := db.MongoClient.Collection("users")
	ctx := context.Background()

	filter := bson.M{
		"_id": u.ID,
	}

	update := bson.M{
		"$set": u,
	}

	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	doc := coll.FindOneAndUpdate(ctx, filter, update, opts)
	if doc.Err() != nil {
		if doc.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Printf("error updating user %s: %s \n", u.ID, doc.Err())
		return nil, doc.Err()
	}

	var user User
	if err := doc.Decode(&user); err != nil {
		log.Printf("error decoding user %s: %s", u.ID, err)
		return nil, err
	}

	return &user, nil
}
