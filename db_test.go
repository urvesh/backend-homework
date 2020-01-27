package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	os.Setenv("MONGO_HOST", "localhost")
	os.Setenv("MONGO_PORT", "27017")
	assert.NotNil(t, NewDB(), "db object should not be nil")
}

func Test_createUserData(t *testing.T) {
	users := createUserData()
	assert.NotNil(t, users, "users array should not be nil")

	// make sure test data contain one of these user Ids
	ids := map[string]bool{
		"5e2e39ee290f5a56ffda9ed5": true,
		"5e2e39ee290f5a56ffda9ed6": true,
		"5e2e39ee290f5a56ffda9ed7": true,
		"5e2e39ee290f5a56ffda9ed8": true,
		"5e2e39ee290f5a56ffda9ed9": true,
		"5e2e39ee290f5a56ffda9eda": true,
	}

	for _, v := range users {
		if _, ok := ids[v.ID]; !ok {
			t.Errorf("user id %s does not belong in known userIds list", v.ID)
		}
	}
}

func Test_createRatingsData(t *testing.T) {
	likes := createRatingsData()
	assert.NotNil(t, likes, "ratings data is not nil")

	// make sure test data only contain these user ids
	ids := map[string]bool{
		"5e2e39ee290f5a56ffda9ed5": true,
		"5e2e39ee290f5a56ffda9ed6": true,
		"5e2e39ee290f5a56ffda9ed7": true,
		"5e2e39ee290f5a56ffda9ed8": true,
		"5e2e39ee290f5a56ffda9ed9": true,
		"5e2e39ee290f5a56ffda9eda": true,
	}

	for _, v := range likes {
		if _, ok := ids[v.FromUserID]; !ok {
			t.Errorf("fromUserId %s does not belong in known userIds list", v.FromUserID)
		}
		if _, ok := ids[v.ToUserID]; !ok {
			t.Errorf("toUserId %s does not belong in known userIds list", v.ToUserID)
		}
	}
}
