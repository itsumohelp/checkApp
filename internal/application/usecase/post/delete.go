package post

import (
	"context"

	"checkapp/internal/domain/repository"
)

type DeletePostInput struct {
	ID string
}

type DeletePostUseCase interface {
	Execute(ctx context.Context, input DeletePostInput) error
}

type deletePostInteractor struct {
	postRepo repository.PostRepository
}

func NewDeletePostUseCase(postRepo repository.PostRepository) DeletePostUseCase {
	return &deletePostInteractor{postRepo: postRepo}
}

func (uc *deletePostInteractor) Execute(ctx context.Context, input DeletePostInput) error {
	if _, err := uc.postRepo.FindByID(ctx, input.ID); err != nil {
		return err
	}
	return uc.postRepo.Delete(ctx, input.ID)
}
