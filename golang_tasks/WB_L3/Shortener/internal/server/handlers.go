package server

import (
	"context"
	"net/http"

	"shortener/internal/shortener"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *shortener.Service
}

func NewHandlers(service *shortener.Service) *Handlers {
	return &Handlers{service: service}
}

// POST /shorten
func (h *Handlers) Short(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	newURL, err := h.service.Shorten(context.Background(), req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": newURL})
}

// GET /s/:code
func (h *Handlers) Redirect(c *gin.Context) {
	code := c.Param("code")

	orig, err := h.service.Resolve(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	h.service.RecordVisit(code, c.Request.UserAgent())
	c.Redirect(http.StatusFound, orig)
}

// GET /analytics/:code
func (h *Handlers) Analytics(c *gin.Context) {
	code := c.Param("code")

	visits, err := h.service.GetAnalytics(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no analytics"})
		return
	}

	c.JSON(http.StatusOK, visits)
}
