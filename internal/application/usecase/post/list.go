package post

import (
	"context"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type ListPostsInput struct {
	Status *entity.PostStatus
	Page   int
	Limit  int
}

type ListPostsOutput struct {
	Posts []*entity.Post
	Total int64
	Page  int
	Limit int
}

type ListPostsUseCase interface {
	Execute(ctx context.Context, input ListPostsInput) (*ListPostsOutput, error)
}

type listPostsInteractor struct {
	postRepo repository.PostRepository
}

func NewListPostsUseCase(postRepo repository.PostRepository) ListPostsUseCase {
	return &listPostsInteractor{postRepo: postRepo}
}

func (uc *listPostsInteractor) Execute(ctx context.Context, input ListPostsInput) (*ListPostsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 || input.Limit > 100 {
		input.Limit = 20
	}

	posts, total, err := uc.postRepo.FindAll(ctx, repository.PostFilter{
		Status: input.Status,
		Page:   input.Page,
		Limit:  input.Limit,
	})
	if err != nil {
		return nil, err
	}
	return &ListPostsOutput{Posts: posts, Total: total, Page: input.Page, Limit: input.Limit}, nil
}
