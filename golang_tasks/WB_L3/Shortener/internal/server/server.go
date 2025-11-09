package server

import (
	"shortener/internal/shortener"

	"github.com/gin-gonic/gin"
)

func NewServer(svc *shortener.Service) *gin.Engine {
	router := gin.Default()
	h := NewHandlers(svc)

	router.POST("/shorten", h.Short)
	router.GET("/s/:code", h.Redirect)
	router.GET("/analytics/:code", h.Analytics)

	return router
}
