package post

import (
	"context"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type UpdatePostInput struct {
	ID      string
	Title   string
	Content string
}

type UpdatePostOutput struct {
	Post *entity.Post
}

type UpdatePostUseCase interface {
	Execute(ctx context.Context, input UpdatePostInput) (*UpdatePostOutput, error)
}

type updatePostInteractor struct {
	postRepo repository.PostRepository
}

func NewUpdatePostUseCase(postRepo repository.PostRepository) UpdatePostUseCase {
	return &updatePostInteractor{postRepo: postRepo}
}

func (uc *updatePostInteractor) Execute(ctx context.Context, input UpdatePostInput) (*UpdatePostOutput, error) {
	post, err := uc.postRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := post.Update(input.Title, input.Content); err != nil {
		return nil, err
	}
	if err := uc.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}
	return &UpdatePostOutput{Post: post}, nil
}
