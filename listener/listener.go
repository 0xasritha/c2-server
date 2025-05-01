package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Listener interface {
	listen()
}

func NewListener(db *gorm.DB, commProtocol CommProtocol) (Listener, error) {
	switch commProtocol {
	case HTTPS:
		return &HTTPSListener{
			db: db,
		}, nil
	case DNS:
		return &DNSListener{}, nil
	}
	return nil, fmt.Errorf("invalid communication protocol: %s", commProtocol)
}

type HTTPSListener struct {
	db *gorm.DB
}

func (l *HTTPSListener) listen() {
	router := gin.Default()

	router.RunTLS(":443", "./tls.cert", "./tls.key") // TODO: cert + key paths should not be hardcoded

	implant := router.Group("/implant")
	implant.POST("/register", l.Register)
	implant.POST("/beacon", l.Beacon)

	operator := router.Group("/operator")
	operator.GET("/implants", l.ListImplants)
	operator.GET("/implant/:id", l.GetImplant)
	operator.GET("/tasks/", l.ListTasks)
	operator.GET("/tasks/:task-id", l.GetTask)
	operator.POST("/tasks/:implant-id", l.AddTask)
	operator.GET("/results/:implant-id", l.ListResults)
}

func (l *HTTPSListener) Register(c *gin.Context) {
	var newImplant *Implant
	if err := c.ShouldBind(&newImplant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func (l *HTTPSListener) Beacon(c *gin.Context) {}

func (l *HTTPSListener) ListImplants(c *gin.Context) {

}

func (l *HTTPSListener) GetImplant(c *gin.Context) {}

func (l *HTTPSListener) ListTasks(c *gin.Context)   {}
func (l *HTTPSListener) GetTask(c *gin.Context)     {}
func (l *HTTPSListener) AddTask(c *gin.Context)     {}
func (l *HTTPSListener) ListResults(c *gin.Context) {}

type DNSListener struct {
}

func (listener *DNSListener) listen() {}
