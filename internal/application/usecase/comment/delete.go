package comment

import (
	"context"

	"checkapp/internal/domain/repository"
)

type DeleteCommentInput struct {
	CommentID string
	PostID    string
}

type DeleteCommentUseCase interface {
	Execute(ctx context.Context, input DeleteCommentInput) error
}

type deleteCommentInteractor struct {
	commentRepo repository.CommentRepository
}

func NewDeleteCommentUseCase(commentRepo repository.CommentRepository) DeleteCommentUseCase {
	return &deleteCommentInteractor{commentRepo: commentRepo}
}

func (uc *deleteCommentInteractor) Execute(ctx context.Context, input DeleteCommentInput) error {
	comment, err := uc.commentRepo.FindByID(ctx, input.CommentID)
	if err != nil {
		return err
	}
	if comment.PostID() != input.PostID {
		return repository.ErrCommentNotFound
	}
	return uc.commentRepo.Delete(ctx, input.CommentID)
}
