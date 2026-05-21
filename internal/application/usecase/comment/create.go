package comment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type CreateCommentInput struct {
	PostID  string
	Author  string
	Content string
}

type CreateCommentOutput struct {
	Comment *entity.Comment
}

type CreateCommentUseCase interface {
	Execute(ctx context.Context, input CreateCommentInput) (*CreateCommentOutput, error)
}

type createCommentInteractor struct {
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
}

func NewCreateCommentUseCase(postRepo repository.PostRepository, commentRepo repository.CommentRepository) CreateCommentUseCase {
	return &createCommentInteractor{postRepo: postRepo, commentRepo: commentRepo}
}

func (uc *createCommentInteractor) Execute(ctx context.Context, input CreateCommentInput) (*CreateCommentOutput, error) {
	post, err := uc.postRepo.FindByID(ctx, input.PostID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrPostNotFound
		}
		return nil, err
	}
	if post.Status() != entity.PostStatusPublished {
		return nil, entity.ErrCommentPostNotPublished
	}

	comment, err := entity.NewComment(uuid.New().String(), input.PostID, input.Author, input.Content)
	if err != nil {
		return nil, err
	}
	if err := uc.commentRepo.Save(ctx, comment); err != nil {
		return nil, err
	}
	return &CreateCommentOutput{Comment: comment}, nil
}
