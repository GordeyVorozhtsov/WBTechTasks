package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"salestracker/internal/tracker"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *tracker.Service
}

func NewHandlers(service *tracker.Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) GetItems(c *gin.Context) {
	items, err := h.service.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handlers) PostItems(c *gin.Context) {
	var i tracker.Item
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if i.OccurredAt.IsZero() {
		i.OccurredAt = time.Now()
	}
	if err := h.service.Create(context.Background(), &i); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handlers) DeleteItems(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handlers) PutItems(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var i tracker.Item
	if err := c.ShouldBindJSON(&i); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	i.ID = id
	if err := h.service.Update(context.Background(), &i); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GET /analytics?from=2025-01-01&to=2025-12-31
func (h *Handlers) GetAnalytics(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	var err error

	if fromStr == "" {
		from = time.Now().AddDate(0, -1, 0) // по умолчанию за последний месяц
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date"})
			return
		}
	}

	if toStr == "" {
		to = time.Now()
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date"})
			return
		}
	}

	res, err := h.service.Analytics(context.Background(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handlers) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}
