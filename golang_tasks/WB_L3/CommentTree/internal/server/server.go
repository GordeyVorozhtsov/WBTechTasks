package server

import (
	"commenttree/internal/commenter"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewServer(svc *commenter.Service) http.Handler {
	router := gin.Default()
	h := NewHandlers(svc)

	router.LoadHTMLGlob("./templates/*")
	router.GET("/", h.IndexPage)

	router.POST("/comments", h.PostComment)      // создать комментарий
	router.GET("/comments", h.GetComment)        // получить дерево или топовые комментарии
	router.DELETE("/comments/:id", h.DelComment) // удалить комментарий с вложенными

	return router
}
