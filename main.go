// package main is a small API that allows users to see and interact with each other
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// appContext holds application level config
type appContext struct {
	DB *DB
}

func main() {
	// .env file might not exist, but envars might..
	if err := godotenv.Load(); err != nil {
		log.Println("warning no .env file found.")
	}

	// setup app config
	app := appContext{
		DB: NewDB(),
	}

	// load default data in database
	PopulateDatabase(app.DB)

	// register handlers
	r := gin.Default()
	r.GET("/users", app.getAllUsers)
	r.GET("/users/:id/likes", app.getIncomingLikes)
	r.PUT("/users/:id", app.editUser)
	r.POST("/users/:id/ratings", app.newRating)
	r.GET("/users/:id/matches", app.getMatches)

	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal("error running server ", err)
	}
}

// see all users that exist within db. helps to get user ids for testing
func (app *appContext) getAllUsers(c *gin.Context) {
	userIDs, err := FindAllUsers(app.DB)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userIDs,
	})
	return
}

// get all incoming likes for a particular user
func (app *appContext) getIncomingLikes(c *gin.Context) {
	userId := c.Param("id")

	users, err := FindIncomingLikes(app.DB, userId)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
	return
}

// edit a user
func (app *appContext) editUser(c *gin.Context) {
	userId := c.Param("id")

	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		if err == io.EOF {
			errorResponse(c, http.StatusBadRequest, NewErrorf("invalid request body: %s", err))
			return
		}
		errorResponse(c, http.StatusInternalServerError, NewErrorf("error binding to user struct: %s", err))
		return
	}

	// user can only change certain number of fields depending on contract.
	u.ID = userId

	// zero out time so it cannot be overridden
	u.CreatedDate = time.Time{}

	user, err := u.Edit(app.DB)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		errorResponse(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
	return
}

// add a new like entry. Does not validate if user ids exist within the system
func (app *appContext) newRating(c *gin.Context) {
	id := c.Param("id")
	var r Rating

	if err := c.ShouldBindJSON(&r); err != nil {
		errorResponse(c, http.StatusInternalServerError, NewErrorf("error binding data to Rating object %s", err))
		return
	}

	if r.ToUserID == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("missing one or more required fields: toUserId, type"))
		return
	}

	if r.Type != LIKE && r.Type != BLOCK && r.Type != REPORT {
		errorResponse(c, http.StatusBadRequest, errors.New(fmt.Sprintf("type must be either %s, %s, or %s", LIKE, BLOCK, REPORT)))
		return
	}

	if r.Type == REPORT && r.Reason == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("reason cannot be blank"))
		return
	}

	r.FromUserID = id

	if err := r.Save(app.DB); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
	return
}

// gets users who have been matched up to this userId
func (app *appContext) getMatches(c *gin.Context) {
	id := c.Param("id")

	users, err := FindMatches(app.DB, id)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
	return
}

// helper function to return 500 errors
func errorResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

// helper function to print json into console
func prettyPrint(i interface{}) {
	b, _ := json.MarshalIndent(i, "", "    ")
	fmt.Println(string(b))
}

// NewError will log and return a new instance of error
func NewErrorf(format string, v ...interface{}) error {
	msg := fmt.Sprintf(format, v...)
	log.Println(msg)
	return errors.New(msg)
}
