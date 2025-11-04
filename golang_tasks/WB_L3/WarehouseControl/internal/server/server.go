package server

import (
	"net/http"
	"os"
	"strings"

	controller "warehousecontrol/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

const jwtSecretEnv = "JWT_SECRET"

// Middleware проверка JWT и складывание роли в контекст
func authMiddleware() gin.HandlerFunc {
	secret := []byte(strings.TrimSpace(getenv(jwtSecretEnv, "supersecret")))
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}
		tokenStr := parts[1]
		tok, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := tok.Claims.(*jwtClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		// кладём роль в контекст
		c.Set("role", claims.Role)
		c.Next()
	}
}

func getenv(k, def string) string {
	if v := strings.TrimSpace(strings.Trim(os.Getenv(k), " ")); v != "" {
		return v
	}
	return def
}

func NewServer(svc *controller.Service) *gin.Engine {
	router := gin.Default()
	h := NewHandlers(svc)

	router.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	router.POST("/login", h.Login)

	apiv := router.Group("/api")
	apiv.Use(authMiddleware())
	{
		apiv.GET("/items", h.GetItems)
		apiv.POST("/items", h.PostItem)
		apiv.PUT("/items/:id", h.PutItem)
		apiv.DELETE("/items/:id", h.DeleteItem)
		apiv.GET("/items/:id/history", h.GetItemHistory)
	}

	router.Static("/static", "./static")
	return router
}
