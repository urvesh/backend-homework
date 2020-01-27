package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initAppContext() *appContext {
	return &appContext{
		DB: NewDB(),
	}
}

func TestFindAllUsers(t *testing.T) {
	app := initAppContext()
	router := setupRouter(app)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var m map[string][]*User
	err := json.Unmarshal(w.Body.Bytes(), &m)
	assert.Nil(t, err)

	assert.NotNil(t, m["data"])
}

func TestNewErrorf(t *testing.T) {
	want := errors.New("test 1")
	got := NewErrorf("test %s", "1")
	assert.Equal(t, want, got)
}
