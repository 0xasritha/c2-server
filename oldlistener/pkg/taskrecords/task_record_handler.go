package taskrecords

import (
	"asritha.dev/listener/pkg/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Store             store.MemStore[TaskRecord]
	taskRecordCounter uint32
}

func NewHandler() *Handler {
	return &Handler{
		taskRecordCounter: 0,
	}
}

func (h *Handler) ListTaskRecords(c *gin.Context) {
	taskRecords, _ := h.Store.List() // TODO: handle error

	if len(taskRecords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no results found"})
	}

	c.JSON(http.StatusOK, gin.H{"taskRecords": taskRecords})
}
