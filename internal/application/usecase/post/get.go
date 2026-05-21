package post

import (
	"context"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type GetPostInput struct {
	ID string
}

type GetPostOutput struct {
	Post *entity.Post
}

type GetPostUseCase interface {
	Execute(ctx context.Context, input GetPostInput) (*GetPostOutput, error)
}

type getPostInteractor struct {
	postRepo repository.PostRepository
}

func NewGetPostUseCase(postRepo repository.PostRepository) GetPostUseCase {
	return &getPostInteractor{postRepo: postRepo}
}

func (uc *getPostInteractor) Execute(ctx context.Context, input GetPostInput) (*GetPostOutput, error) {
	post, err := uc.postRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return &GetPostOutput{Post: post}, nil
}
