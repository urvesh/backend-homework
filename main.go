package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

/**
Seeing and interacting with recommendations of potential matches
Seeing and interacting with others who have already liked you
Seeing and interacting with existing matches
Profile and account editing
*/

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
	r.POST("/likes", app.newLike)

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

// get incoming user likes for a particular user id
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
		log.Printf("error reading user edit data: %s\n", err)
		errorResponse(c, http.StatusInternalServerError, err)
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
func (app *appContext) newLike(c *gin.Context) {
	var l Like

	if err := c.ShouldBindJSON(&l); err != nil {
		log.Printf("error binding data to Like object %s\n", err)
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if l.UserID == "" || l.LikeUserID == "" {
		errorResponse(c, http.StatusBadRequest, errors.New("missing one or more required fields: userId, likeUserId"))
		return
	}

	if err := l.Save(app.DB); err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, nil)
	return
}

// helper function to return 500 errors
func errorResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error": err.Error(),
	})
}

func prettyPrint(i interface{}) {
	b, _ := json.MarshalIndent(i, "", "    ")
	fmt.Println(string(b))
}
