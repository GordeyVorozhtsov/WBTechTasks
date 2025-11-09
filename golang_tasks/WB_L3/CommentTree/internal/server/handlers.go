package server

import (
	"net/http"
	"strconv"

	"commenttree/internal/commenter"

	"github.com/gin-gonic/gin"
)

type createReq struct {
	ParentID *int64 `json:"parent_id"`
	Content  string `json:"content"`
}

type Handlers struct {
	service *commenter.Service
}

func NewHandlers(service *commenter.Service) *Handlers {
	return &Handlers{service: service}
}

// POST /comments
func (h *Handlers) PostComment(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	comment := h.service.Create(req.ParentID, req.Content)
	c.JSON(http.StatusCreated, comment)
}

// GET /comments?parent={id}
func (h *Handlers) GetComment(c *gin.Context) {
	parent := c.Query("parent")
	if parent == "" {
		page := atoiWithDefault(c.Query("page"), 1)
		size := atoiWithDefault(c.Query("page_size"), 10)
		sortBy := c.DefaultQuery("sort", "created_desc")

		comments, total := h.service.ListTop(page, size, sortBy)
		c.JSON(http.StatusOK, gin.H{"items": comments, "total": total})
		return
	}

	id, err := strconv.ParseInt(parent, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent id"})
		return
	}

	tree, err := h.service.GetTree(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}
	c.JSON(http.StatusOK, tree)
}

// DELETE /comments/:id
func (h *Handlers) DelComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func atoiWithDefault(sval string, def int) int {
	if sval == "" {
		return def
	}
	v, err := strconv.Atoi(sval)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

func (h *Handlers) IndexPage(c *gin.Context) {
	query := c.Query("q")

	var comments []*commenter.Comment
	var total int

	if query != "" {
		comments, total = h.service.Search(query, 1, 100)
	} else {
		comments, total = h.service.ListTop(1, 100, "created_desc")
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"comments": comments,
		"total":    total,
		"query":    query,
	})
}
