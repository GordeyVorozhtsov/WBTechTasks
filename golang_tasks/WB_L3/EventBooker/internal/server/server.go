package server

import (
	"eventbooker/internal/booker"

	"github.com/gin-gonic/gin"
)

func NewServer(svc *booker.Service) *gin.Engine {
	router := gin.Default()
	h := NewHandlers(svc)

	router.POST("/events", h.CreateEvent)
	router.POST("/events/:id/book", h.BookEvent)
	router.POST("/events/:id/confirm", h.ConfirmEvent)
	router.GET("/events/:id", h.GetEvent)
	router.GET("/events", h.GetAllEvents)

	router.LoadHTMLGlob("templates/*")
	router.GET("/", h.IndexPage)

	return router
}
