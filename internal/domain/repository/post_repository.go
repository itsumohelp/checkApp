package repository

import (
	"context"
	"errors"

	"checkapp/internal/domain/entity"
)

var ErrPostNotFound = errors.New("post not found")

type PostFilter struct {
	Status *entity.PostStatus
	Page   int
	Limit  int
}

type PostRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Post, error)
	FindAll(ctx context.Context, filter PostFilter) ([]*entity.Post, int64, error)
	Save(ctx context.Context, post *entity.Post) error
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id string) error
}
