package server

import (
	kafka "github.com/wb-go/wbf/kafka"

	"github.com/gin-gonic/gin"
)

func NewServer(prod *kafka.Producer) *gin.Engine {
	router := gin.Default()
	h := NewHandlers(prod)

	router.Static("/storage", "./storage/processed")
	router.StaticFile("/", "./web/index.html")

	router.POST("/upload", h.UploadImg)
	router.GET("/image/:id", h.GetImg)
	router.DELETE("/image/:id", h.DelImg)

	return router
}
