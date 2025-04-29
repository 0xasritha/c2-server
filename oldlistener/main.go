package main

import (
	"asritha.dev/oldlistener/pkg/results"
	"asritha.dev/oldlistener/pkg/taskrecords"
	"asritha.dev/oldlistener/pkg/tasks"
	"github.com/gin-gonic/gin"
)

var (
	taskRecordsHandler = taskrecords.NewHandler()
	tasksHandler       = tasks.NewHandler(taskRecordsHandler)
	listener           = Listener{
		resultsHandler: results.NewHandler(tasksHandler, taskRecordsHandler),
		tasksHandler:   tasksHandler,
	}
)

func main() {
	router := gin.Default()

	// `TASKS` endpoints
	router.GET("/tasks", listener.tasksHandler.ListTasks)
	router.POST("/tasks", listener.tasksHandler.AddTasks)

	// `RESULTS` ENDPOINTS
	router.GET("/results", listener.resultsHandler.ListResults)
	router.GET("/results", listener.resultsHandler.AddResults)

	// `TASK RECORD` ENDPOINTS
	// TODO: what is naming convention for API paths? taskRecord or taskrecords?
	router.GET("/taskrecords", listener.taskRecordsHandler.ListTaskRecords)

	router.Run(":8080")
}
