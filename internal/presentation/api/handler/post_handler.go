package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	postUC "checkapp/internal/application/usecase/post"
	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type PostHandler struct {
	createPost    postUC.CreatePostUseCase
	getPost       postUC.GetPostUseCase
	listPosts     postUC.ListPostsUseCase
	updatePost    postUC.UpdatePostUseCase
	deletePost    postUC.DeletePostUseCase
	publishPost   postUC.PublishPostUseCase
	unpublishPost postUC.UnpublishPostUseCase
}

func NewPostHandler(
	create postUC.CreatePostUseCase,
	get postUC.GetPostUseCase,
	list postUC.ListPostsUseCase,
	update postUC.UpdatePostUseCase,
	delete postUC.DeletePostUseCase,
	publish postUC.PublishPostUseCase,
	unpublish postUC.UnpublishPostUseCase,
) *PostHandler {
	return &PostHandler{
		createPost:    create,
		getPost:       get,
		listPosts:     list,
		updatePost:    update,
		deletePost:    delete,
		publishPost:   publish,
		unpublishPost: unpublish,
	}
}

type createPostRequest struct {
	Title   string `json:"title"   binding:"required"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author"  binding:"required"`
}

type updatePostRequest struct {
	Title   string `json:"title"   binding:"required"`
	Content string `json:"content" binding:"required"`
}

type postResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type listResponse struct {
	Posts []postResponse `json:"posts"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func toResponse(p *entity.Post) postResponse {
	return postResponse{
		ID:        p.ID(),
		Title:     p.Title(),
		Content:   p.Content(),
		Author:    p.Author(),
		Status:    string(p.Status()),
		CreatedAt: p.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: p.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, entity.ErrInvalidTitle),
		errors.Is(err, entity.ErrInvalidContent),
		errors.Is(err, entity.ErrInvalidAuthor):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, entity.ErrAlreadyPublished),
		errors.Is(err, entity.ErrAlreadyDraft):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.createPost.Execute(c.Request.Context(), postUC.CreatePostInput{
		Title: req.Title, Content: req.Content, Author: req.Author,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toResponse(out.Post))
}

func (h *PostHandler) GetByID(c *gin.Context) {
	out, err := h.getPost.Execute(c.Request.Context(), postUC.GetPostInput{ID: c.Param("id")})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(out.Post))
}

func (h *PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	input := postUC.ListPostsInput{Page: page, Limit: limit}
	if s := c.Query("status"); s != "" {
		status := entity.PostStatus(s)
		input.Status = &status
	}

	out, err := h.listPosts.Execute(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	posts := make([]postResponse, len(out.Posts))
	for i, p := range out.Posts {
		posts[i] = toResponse(p)
	}
	c.JSON(http.StatusOK, listResponse{Posts: posts, Total: out.Total, Page: out.Page, Limit: out.Limit})
}

func (h *PostHandler) Update(c *gin.Context) {
	var req updatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.updatePost.Execute(c.Request.Context(), postUC.UpdatePostInput{
		ID: c.Param("id"), Title: req.Title, Content: req.Content,
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(out.Post))
}

func (h *PostHandler) Delete(c *gin.Context) {
	if err := h.deletePost.Execute(c.Request.Context(), postUC.DeletePostInput{ID: c.Param("id")}); err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PostHandler) Publish(c *gin.Context) {
	out, err := h.publishPost.Execute(c.Request.Context(), postUC.PublishPostInput{ID: c.Param("id")})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(out.Post))
}

func (h *PostHandler) Unpublish(c *gin.Context) {
	out, err := h.unpublishPost.Execute(c.Request.Context(), postUC.UnpublishPostInput{ID: c.Param("id")})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponse(out.Post))
}
