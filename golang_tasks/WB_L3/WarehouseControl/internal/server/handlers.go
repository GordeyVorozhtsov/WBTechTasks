package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	controller "warehousecontrol/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handlers struct {
	svc *controller.Service
}

func NewHandlers(s *controller.Service) *Handlers { return &Handlers{svc: s} }

func (h *Handlers) Login(c *gin.Context) {
	var in struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	role := strings.ToLower(in.Role)
	if role != "admin" && role != "manager" && role != "viewer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown role"})
		return
	}
	secret := []byte(getenv(jwtSecretEnv, "supersecret"))
	claims := jwtClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": ts, "role": role})
}

func getRole(c *gin.Context) string {
	if v, ok := c.Get("role"); ok {
		if s, ok2 := v.(string); ok2 {
			return s
		}
	}
	return ""
}

// Authorization helpers
func allowRole(role string, allowed ...string) bool {
	for _, a := range allowed {
		if a == role {
			return true
		}
	}
	return false
}

// GET /api/items
func (h *Handlers) GetItems(c *gin.Context) {
	items, err := h.svc.ListItems(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /api/items
func (h *Handlers) PostItem(c *gin.Context) {
	role := getRole(c)
	if !allowRole(role, "admin", "manager") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var in controller.Item
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	// set defaults
	if in.SKU == "" || in.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sku and name required"})
		return
	}
	ctx := context.Background()
	if err := h.svc.CreateItem(ctx, role, &in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, in)
}

// PUT /api/items/:id
func (h *Handlers) PutItem(c *gin.Context) {
	role := getRole(c)
	if !allowRole(role, "admin", "manager") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var in controller.Item
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	in.ID = id
	ctx := context.Background()
	if err := h.svc.UpdateItem(ctx, role, &in); err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, in)
}

// DELETE /api/items/:id
func (h *Handlers) DeleteItem(c *gin.Context) {
	role := getRole(c)
	if !allowRole(role, "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.svc.DeleteItem(context.Background(), role, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /api/items/:id/history
func (h *Handlers) GetItemHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	rows, err := h.svc.ListAuditForItem(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}
