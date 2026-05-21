package comment

import (
	"context"
	"time"

	"gorm.io/gorm"

	"checkapp/internal/domain/entity"
)

type ListCommentsInput struct {
	PostID string
	Page   int
	Limit  int
}

type ListCommentsOutput struct {
	Comments []*entity.Comment
	Total    int64
	Page     int
	Limit    int
}

type ListCommentsUseCase interface {
	Execute(ctx context.Context, input ListCommentsInput) (*ListCommentsOutput, error)
}

type commentRow struct {
	ID        string
	PostID    string
	Author    string
	Content   string
	CreatedAt time.Time
}

type listCommentsInteractor struct {
	db *gorm.DB
}

func NewListCommentsUseCase(db *gorm.DB) ListCommentsUseCase {
	return &listCommentsInteractor{db: db}
}

func (uc *listCommentsInteractor) Execute(ctx context.Context, input ListCommentsInput) (*ListCommentsOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 || input.Limit > 50 {
		input.Limit = 20
	}

	var total int64
	uc.db.WithContext(ctx).Table("comments").Where("post_id = ?", input.PostID).Count(&total)

	var rows []commentRow
	offset := (input.Page - 1) * input.Limit
	uc.db.WithContext(ctx).Table("comments").
		Where("post_id = ?", input.PostID).
		Order("created_at DESC").
		Offset(offset).Limit(input.Limit).
		Find(&rows)

	comments := make([]*entity.Comment, len(rows))
	for i, r := range rows {
		comments[i] = entity.ReconstituteComment(r.ID, r.PostID, r.Author, r.Content, r.CreatedAt)
	}
	return &ListCommentsOutput{Comments: comments, Total: total, Page: input.Page, Limit: input.Limit}, nil
}
