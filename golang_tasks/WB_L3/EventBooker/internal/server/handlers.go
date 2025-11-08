package server

import (
	"eventbooker/internal/booker"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *booker.Service
}

func NewHandlers(service *booker.Service) *Handlers {
	return &Handlers{service: service}
}

type CreateEventRequest struct {
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Capacity int       `json:"capacity"`
	Timeout  int       `json:"timeout"`
}

type BookEventRequest struct {
	UserEmail string `json:"user_email"`
	Places    int    `json:"places"`
}

type ConfirmBookingRequest struct {
	UserEmail string `json:"user_email"`
}

// POST /events — создание мероприятия;
func (h *Handlers) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(c, req.Title, req.Date, req.Capacity, req.Timeout)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, event)
}

// POST /events/{id}/book
func (h *Handlers) BookEvent(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))
	var req BookEventRequest
	c.ShouldBindJSON(&req)

	booking, err := h.service.BookEvent(c, eventID, req.UserEmail, req.Places)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, booking)
}

// POST /events/{id}/confirm
func (h *Handlers) ConfirmEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid event id"})
		return
	}

	var req ConfirmBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.UserEmail == "" {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	booking, err := h.service.ConfirmBooking(c, eventID, req.UserEmail)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, booking)
}

// GET /events/{id}
func (h *Handlers) GetEvent(c *gin.Context) {
	eventID, _ := strconv.Atoi(c.Param("id"))
	event, err := h.service.GetEvent(c, eventID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, event)
}

func (h *Handlers) GetAllEvents(c *gin.Context) {
	events, err := h.service.GetAllEvents(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, events)
}

func (h *Handlers) IndexPage(c *gin.Context) {
	events, _ := h.service.GetAllEvents(c)
	c.HTML(200, "index.html", gin.H{"events": events})
}
