package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User holds information related to the user collection
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
		err = NewErrorf("error finding users from mongo: %s", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var u User
		if err := cur.Decode(&u); err != nil {
			err = NewErrorf("error decoding into user struct: %s", err)
			return nil, err
		}

		users = append(users, &u)
	}

	if err := cur.Err(); err != nil {
		err = NewErrorf("mongo error: %s", err)
		return nil, err
	}

	return users, nil
}

// FindUserById lookup user by id
func FindUserByID(db *DB, id string) (*User, error) {
	coll := db.MongoClient.Collection("users")

	filter := bson.M{
		"_id": id,
	}

	doc := coll.FindOne(context.Background(), filter)
	if doc.Err() != nil {
		if doc.Err() == mongo.ErrNoDocuments {
			// document not found, not an actual error
			return nil, nil
		}
		return nil, NewErrorf("error looking up user %s: %s", id, doc.Err())
	}

	var u User
	if err := doc.Decode(&u); err != nil {
		return nil, NewErrorf("error decoding user %s into struct: %s", id, err)
	}

	return &u, nil
}

// FindIncomingLikes finds all the users who have liked the given userId
func FindIncomingLikes(db *DB, userId string) ([]*User, error) {
	// find likes where toUserId is this user
	p := &RatingParams{
		Filter: &Rating{
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

// FindMatches gets all the matches this user has
func FindMatches(db *DB, userId string) ([]*User, error) {
	// look up everyone this user likes
	p := &RatingParams{
		Filter: &Rating{
			FromUserID: userId,
			Type:       LIKE,
		},
	}

	myLikes, err := FindRatings(db, p)
	if err != nil {
		return nil, err
	}

	// find who likes this user
	p.Filter.FromUserID = ""
	p.Filter.ToUserID = userId

	incomingLikes, err := FindRatings(db, p)
	if err != nil {
		return nil, err
	}

	// hold ids of outgoing likes
	m := make(map[string]bool)
	for _, v := range myLikes {
		m[v.ToUserID] = true
	}

	matches := make([]*User, 0)

	// loop through incoming likes, find common ids against outgoing likes and look up the user
	for _, v := range incomingLikes {
		if _, ok := m[v.FromUserID]; ok {
			u, err := FindUserByID(db, v.FromUserID)
			if err != nil {
				return nil, err
			}
			matches = append(matches, u)
		}
	}

	return matches, nil
}

// Edit overrides user data with the incoming values
func (u *User) Edit(db *DB) (*User, error) {
	coll := db.MongoClient.Collection("users")
	ctx := context.Background()

	filter := bson.M{
		"_id": u.ID,
	}

	// only set fields that contain a value
	update := bson.M{
		"$set": u,
	}

	// update and return the updated document
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	doc := coll.FindOneAndUpdate(ctx, filter, update, opts)
	if doc.Err() != nil {
		if doc.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, NewErrorf("error updating user %s: %s", u.ID, doc.Err())
	}

	var user User
	if err := doc.Decode(&user); err != nil {
		return nil, NewErrorf("error decoding user %s: %s", u.ID, err)
	}

	return &user, nil
}
