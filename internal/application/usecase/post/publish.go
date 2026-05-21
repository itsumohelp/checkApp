package post

import (
	"context"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type PublishPostInput struct {
	ID string
}

type PublishPostOutput struct {
	Post *entity.Post
}

type PublishPostUseCase interface {
	Execute(ctx context.Context, input PublishPostInput) (*PublishPostOutput, error)
}

type publishPostInteractor struct {
	postRepo repository.PostRepository
}

func NewPublishPostUseCase(postRepo repository.PostRepository) PublishPostUseCase {
	return &publishPostInteractor{postRepo: postRepo}
}

func (uc *publishPostInteractor) Execute(ctx context.Context, input PublishPostInput) (*PublishPostOutput, error) {
	post, err := uc.postRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := post.Publish(); err != nil {
		return nil, err
	}
	if err := uc.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}
	return &PublishPostOutput{Post: post}, nil
}

type UnpublishPostInput struct {
	ID string
}

type UnpublishPostOutput struct {
	Post *entity.Post
}

type UnpublishPostUseCase interface {
	Execute(ctx context.Context, input UnpublishPostInput) (*UnpublishPostOutput, error)
}

type unpublishPostInteractor struct {
	postRepo repository.PostRepository
}

func NewUnpublishPostUseCase(postRepo repository.PostRepository) UnpublishPostUseCase {
	return &unpublishPostInteractor{postRepo: postRepo}
}

func (uc *unpublishPostInteractor) Execute(ctx context.Context, input UnpublishPostInput) (*UnpublishPostOutput, error) {
	post, err := uc.postRepo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err := post.Unpublish(); err != nil {
		return nil, err
	}
	if err := uc.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}
	return &UnpublishPostOutput{Post: post}, nil
}
