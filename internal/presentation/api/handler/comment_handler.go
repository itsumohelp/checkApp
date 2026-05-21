package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	commentUC "checkapp/internal/application/usecase/comment"
	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

const (
	maxCommentLength = 500
	minCommentWords  = 2
)

var bannedPhrases = []string{"buy now", "click here", "limited offer"}

type CommentHandler struct {
	createComment commentUC.CreateCommentUseCase
	listComments  commentUC.ListCommentsUseCase
	deleteComment commentUC.DeleteCommentUseCase
}

func NewCommentHandler(
	create commentUC.CreateCommentUseCase,
	list commentUC.ListCommentsUseCase,
	delete commentUC.DeleteCommentUseCase,
) *CommentHandler {
	return &CommentHandler{createComment: create, listComments: list, deleteComment: delete}
}

type createCommentRequest struct {
	Author  string `json:"author"  binding:"required"`
	Content string `json:"content" binding:"required"`
}

type commentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type listCommentsResponse struct {
	Comments []commentResponse `json:"comments"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

func toCommentResponse(c *entity.Comment) commentResponse {
	return commentResponse{
		ID:        c.ID(),
		PostID:    c.PostID(),
		Author:    c.Author(),
		Content:   c.Content(),
		CreatedAt: c.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func handleCommentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrCommentNotFound), errors.Is(err, repository.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, entity.ErrInvalidCommentContent), errors.Is(err, entity.ErrInvalidCommentAuthor):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, entity.ErrCommentPostNotPublished):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	postID := c.Param("id")

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	words := strings.Fields(req.Author)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}
	req.Author = strings.Join(words, " ")

	req.Content = strings.TrimSpace(req.Content)
	if len([]rune(req.Content)) > maxCommentLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("comment must be %d characters or less", maxCommentLength)})
		return
	}
	if len(strings.Fields(req.Content)) < minCommentWords {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comment is too short"})
		return
	}
	for _, phrase := range bannedPhrases {
		if strings.Contains(strings.ToLower(req.Content), phrase) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "comment contains prohibited content"})
			return
		}
	}

	out, err := h.createComment.Execute(c.Request.Context(), commentUC.CreateCommentInput{
		PostID:  postID,
		Author:  req.Author,
		Content: req.Content,
	})
	if err != nil {
		handleCommentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toCommentResponse(out.Comment))
}

func (h *CommentHandler) List(c *gin.Context) {
	postID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	out, err := h.listComments.Execute(c.Request.Context(), commentUC.ListCommentsInput{
		PostID: postID, Page: page, Limit: limit,
	})
	if err != nil {
		handleCommentError(c, err)
		return
	}

	responses := make([]commentResponse, len(out.Comments))
	for i, comment := range out.Comments {
		responses[i] = toCommentResponse(comment)
	}
	c.JSON(http.StatusOK, listCommentsResponse{
		Comments: responses, Total: out.Total, Page: out.Page, Limit: out.Limit,
	})
}

func (h *CommentHandler) Delete(c *gin.Context) {
	postID := c.Param("id")
	commentID := c.Param("commentId")

	if err := h.deleteComment.Execute(c.Request.Context(), commentUC.DeleteCommentInput{
		CommentID: commentID, PostID: postID,
	}); err != nil {
		handleCommentError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
