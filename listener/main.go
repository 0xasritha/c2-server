package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// TODO: where to put this?
func setupDB() *gorm.DB {
	// setup DB
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable" // connect w/: `psql -h localhost -p 5432 -U postgres -d postgres`
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// TODO: how to read these logs?
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Result{})
	return db
}

func main() {

	db := setupDB()

	httpsListener, _ := NewListener(db, HTTPS)
	httpsListener.listen()

	dnsListener, _ := NewListener(db, DNS)
	dnsListener.listen()

	/*
		// GET PENDING TASKS (FOR OPERATOR)
		router.GET("/tasks/pending", func(c *gin.Context) {
			// TODO: retrieve all objects with
			pendingTasks := make([]Task, 0)
			result := db.Where("completed = ?", "True").Find(&pendingTasks)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			c.JSON(http.StatusOK, gin.H{"tasks": pendingTasks})
		})

		// POST NEW PENDING TASKS (FOR OPERATOR)
		router.POST("/tasks/pending", func(c *gin.Context) {
			// this will fail, idk if should bind json works for *Task
			var newTasks []*Task // []Task
			if err = c.ShouldBindJSON(&newTasks); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}

			// TODO: check that completed = false, if completed = true, throw error, if Type is valid
			// GORM automatically handles indexes
			result := db.Omit("Results").Create(newTasks) // skip saving "Results" association

			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			}

			c.JSON(http.StatusOK, gin.H{"successfully added tasks": newTasks})
		})

		// GET COMPLETED TASKS (FOR OPERATOR)
		router.GET("/tasks/completed", func(c *gin.Context) {

		})

		// POST RESULTS + GET PENDING TASKS (FOR IMPLANT)
		router.POST("/beacon", func(c *gin.Context) {
			// receive results
			// needs to be []*Result?
			var results []Result
			if err = c.ShouldBindJSON(&results); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // correct error?
			}

			for _, result := range results {
				var correspondingTask Task // better name need
				db.First(&correspondingTask, result.TaskID)
				correspondingTask.Results = append(correspondingTask.Results, result)
				db.Save(&correspondingTask)
			}

		})

	*/
}
