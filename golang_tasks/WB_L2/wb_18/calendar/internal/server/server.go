package server

import (
	"calendar/internal/calendar"

	"github.com/gin-gonic/gin"
)

func NewServer(cal *calendar.Calendar) *gin.Engine {
	router := gin.Default()
	
	// Middleware
	router.Use(LoggingMiddleware())
	
	handlers := NewHandlers(cal)
	
	// Routes
	router.POST("/create_event", handlers.CreateEvent)
	router.POST("/update_event", handlers.UpdateEvent)
	router.POST("/delete_event", handlers.DeleteEvent)
	router.GET("/events_for_day", handlers.EventsForDay)
	router.GET("/events_for_week", handlers.EventsForWeek)
	router.GET("/events_for_month", handlers.EventsForMonth)
	
	return router
}