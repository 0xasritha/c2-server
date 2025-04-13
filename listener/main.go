package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

// TODO: have a models.go w/ the models ?

// TODO: draw a diagram explaining how the different parts of this work together
var (
	listener = KrakenListener{
		TaskManager:   NewLocalTaskStore(),
		ResultManager: NewLocalResultsStore(),
	}
)

func main() {
	router := gin.Default()

	// `TASKS` endpoints
	// `RESULTS` ENDPOINTS
	// get all existing results
	router.GET("/results", func(c *gin.Context) {
		c.JSON(http.StatusOK, listener.getResults())
	})

	// implant can POST results, receive new tasks here
	router.POST("/results", func(c *gin.Context) {
	})

	// `HISTORY` ENDPOINTS
	router.Run(":8080")
}

// TODO: what the resultsCREATE should do:
// check if results returned at empty
// if not empty:
// parse out JSON in request body
// save each Request object to database
// get list of Task objects waiting to be served the implant
// delete Task objects we are serving to the implant (to prevent tasks from being executed twice)
// send Task objects to the implant in the POST request response
// if empty:
// return the Task objects waiting, then delete them
