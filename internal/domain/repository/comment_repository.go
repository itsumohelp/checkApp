package repository

import (
	"context"
	"errors"

	"checkapp/internal/domain/entity"
)

var ErrCommentNotFound = errors.New("comment not found")

type CommentFilter struct {
	PostID string
	Page   int
	Limit  int
}

type CommentRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Comment, error)
	FindByPostID(ctx context.Context, filter CommentFilter) ([]*entity.Comment, int64, error)
	Save(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id string) error
}
