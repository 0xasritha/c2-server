package tasks

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Store             store.MemStore[Task]
	taskRecordHandler *taskrecords.Handler
	taskCounter       uint32 // TODO: this is def not the fastest way to store tasks, map with keys as numbers?
}

func NewHandler(taskRecordHandler *taskrecords.Handler) *Handler {
	return &Handler{
		taskRecordHandler: taskRecordHandler,
		taskCounter:       0,
	}
}

func (h *Handler) ListTasks(c *gin.Context) {
	tasks, err := h.Store.List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) AddTasks(c *gin.Context) {
	var tasks []Task

	// c.ShouldBindJSON converts HTTP body into struct (wrapper around `json.Unmarshal`)
	if err := c.ShouldBindJSON(&tasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, task := range tasks {
		_ = h.Store.Add(h.taskCounter, task) // TODO: handle error
		taskRecord := taskrecords.NewTaskRecord(task.id, task.TaskType, task.Options)
		h.taskRecordHandler.Store.Add(task.id, *taskRecord)
		h.taskCounter++
	}

	c.JSON(http.StatusOK, gin.H{"successfully added tasks": tasks})
}
