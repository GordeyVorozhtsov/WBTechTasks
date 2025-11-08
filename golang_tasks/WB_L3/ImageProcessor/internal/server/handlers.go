package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	kafka "github.com/wb-go/wbf/kafka"
)

type Handlers struct {
	Producer *kafka.Producer
}

func NewHandlers(prod *kafka.Producer) *Handlers {
	return &Handlers{Producer: prod}
}

type ImageTask struct {
	ID  string `json:"id"`
	Ext string `json:"ext"`
}

// POST /upload
func (h *Handlers) UploadImg(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	id := uuid.New().String()
	ext := filepath.Ext(file.Filename)
	filename := "./storage/original/" + id + ext

	os.MkdirAll("./storage/original", os.ModePerm)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
		return
	}

	task := ImageTask{ID: id, Ext: ext}
	data, _ := json.Marshal(task)

	if err := h.Producer.Send(context.Background(), []byte(id), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot enqueue task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "queued"})
}

// GET /image/:id
func (h *Handlers) GetImg(c *gin.Context) {
	id := c.Param("id")
	processedDir := "./storage/processed"

	files, _ := filepath.Glob(filepath.Join(processedDir, id+"*"))
	if len(files) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "processing"})
		return
	}

	c.File(files[0])
}

// DELETE /image/:id
func (h *Handlers) DelImg(c *gin.Context) {
	id := c.Param("id")

	for _, dir := range []string{"original", "processed", "thumbnails"} {
		matches, _ := filepath.Glob(filepath.Join("./storage", dir, id+"*"))
		for _, f := range matches {
			os.Remove(f)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
