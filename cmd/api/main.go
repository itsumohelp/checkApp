package main

import (
	"log"
	"os"

	commentUC "checkapp/internal/application/usecase/comment"
	postUC "checkapp/internal/application/usecase/post"
	"checkapp/internal/infrastructure/database"
	"checkapp/internal/infrastructure/persistence"
	"checkapp/internal/presentation/api"
	"checkapp/internal/presentation/api/handler"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "blog.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.NewSQLite(dbPath)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	postRepo := persistence.NewPostRepository(db)
	commentRepo := persistence.NewCommentRepository(db)

	postHandler := handler.NewPostHandler(
		postUC.NewCreatePostUseCase(postRepo),
		postUC.NewGetPostUseCase(postRepo),
		postUC.NewListPostsUseCase(postRepo),
		postUC.NewUpdatePostUseCase(postRepo),
		postUC.NewDeletePostUseCase(postRepo),
		postUC.NewPublishPostUseCase(postRepo),
		postUC.NewUnpublishPostUseCase(postRepo),
	)

	commentHandler := handler.NewCommentHandler(
		commentUC.NewCreateCommentUseCase(postRepo, commentRepo),
		commentUC.NewListCommentsUseCase(db),
		commentUC.NewDeleteCommentUseCase(commentRepo),
	)

	router := api.NewRouter(postHandler, commentHandler)

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
