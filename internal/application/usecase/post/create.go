package post

import (
	"context"

	"github.com/google/uuid"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type CreatePostInput struct {
	Title   string
	Content string
	Author  string
}

type CreatePostOutput struct {
	Post *entity.Post
}

type CreatePostUseCase interface {
	Execute(ctx context.Context, input CreatePostInput) (*CreatePostOutput, error)
}

type createPostInteractor struct {
	postRepo repository.PostRepository
}

func NewCreatePostUseCase(postRepo repository.PostRepository) CreatePostUseCase {
	return &createPostInteractor{postRepo: postRepo}
}

func (uc *createPostInteractor) Execute(ctx context.Context, input CreatePostInput) (*CreatePostOutput, error) {
	post, err := entity.NewPost(uuid.New().String(), input.Title, input.Content, input.Author)
	if err != nil {
		return nil, err
	}
	if err := uc.postRepo.Save(ctx, post); err != nil {
		return nil, err
	}
	return &CreatePostOutput{Post: post}, nil
}
