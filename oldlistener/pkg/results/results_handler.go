package results

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	Store              store.MemStore[Result]
	tasksHandler       *tasks.Handler
	taskRecordsHandler *taskrecords.Handler
	resultCounter      uint32
}

func NewHandler(tasksHandler *tasks.Handler, taskRecordsHandler *taskrecords.Handler) *Handler {
	return &Handler{
		tasksHandler:       tasksHandler,
		taskRecordsHandler: taskRecordsHandler,
		resultCounter:      0,
	}
}

func (h *Handler) ListResults(c *gin.Context) {
	results, err := h.Store.List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) AddResults(c *gin.Context) {
	var results []Result
	if err := c.ShouldBindJSON(&results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if results in request body are empty; #TODO: test this edge case - does it error out in the above way if it's an empty array ?
	if len(results) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no results found"})
	}

	for _, result := range results {
		h.Store.Add(h.resultCounter, result)

		// TODO: should be getting from map, using key lookup for tasks
		taskRecords, err := h.taskRecordsHandler.Store.List()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// taskRecordID is taskID -> simplify naming?
		// taskRecordID is stored in taskRecord
		for taskRecordID, taskRecord := range taskRecords {
			if result.TaskId == taskRecordID {
				taskRecord.AddResult(result.Content, result.IsSuccess) // will add results to matching task
			}
		}

		h.resultCounter++
	}

	newTasks, err := h.tasksHandler.Store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	// does results need TaskHandler?
	for id, _ := range newTasks {
		h.tasksHandler.Store.Delete(id)
	}

	c.JSON(http.StatusOK, newTasks)
}
