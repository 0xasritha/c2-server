package task

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Store
}

func NewHandler(store Store) *Handler {
	return &Handler{}
}

func (h Handler) CreateTask(c *gin.Context) {

}
func (h Handler) CreateTasks(c *gin.Context) {
	var newTasks []Task
	if err := c.ShouldBindJSON(&newTasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Store.AddTasks(newTasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTasks)
}
func (h Handler) ListTasks(c *gin.Context) {
	tasks, err := h.Store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
func (h Handler) GetTask(c *gin.Context) {
}
func (h Handler) UpdateTask(c *gin.Context) {
}
func (h Handler) DeleteTask(c *gin.Context) {
}
