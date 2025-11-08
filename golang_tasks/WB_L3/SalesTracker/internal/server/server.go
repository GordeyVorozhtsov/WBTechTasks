package server

import (
	"salestracker/internal/tracker"

	"github.com/gin-gonic/gin"
)

func NewServer(svc *tracker.Service) *gin.Engine {
	router := gin.Default()
	h := NewHandlers(svc)

	router.LoadHTMLGlob("templates/*")

	router.GET("/", h.IndexPage)
	router.GET("/items", h.GetItems)
	router.POST("/items", h.PostItems)
	router.PUT("/items/:id", h.PutItems)
	router.DELETE("/items/:id", h.DeleteItems)
	router.GET("/analytics", h.GetAnalytics)

	return router
}
