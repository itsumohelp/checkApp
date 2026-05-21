package api

import (
	"github.com/gin-gonic/gin"

	"checkapp/internal/presentation/api/handler"
	"checkapp/internal/presentation/api/middleware"
)

func NewRouter(postHandler *handler.PostHandler, commentHandler *handler.CommentHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	v1 := r.Group("/api/v1")
	posts := v1.Group("/posts")
	{
		posts.POST("", postHandler.Create)
		posts.GET("", postHandler.List)
		posts.GET("/:id", postHandler.GetByID)
		posts.PUT("/:id", postHandler.Update)
		posts.DELETE("/:id", postHandler.Delete)
		posts.PATCH("/:id/publish", postHandler.Publish)
		posts.PATCH("/:id/unpublish", postHandler.Unpublish)

		posts.POST("/:id/comments", commentHandler.Create)
		posts.GET("/:id/comments", commentHandler.List)
		posts.DELETE("/:id/comments/:commentId", commentHandler.Delete)
	}

	return r
}
