package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type commentModel struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	PostID    string    `gorm:"not null;index"`
	Author    string    `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (commentModel) TableName() string { return "comments" }

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	db.AutoMigrate(&commentModel{}) //nolint:errcheck
	return &commentRepository{db: db}
}

func (r *commentRepository) FindByID(ctx context.Context, id string) (*entity.Comment, error) {
	var m commentModel
	result := r.db.WithContext(ctx).First(&m, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, repository.ErrCommentNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return entity.ReconstituteComment(m.ID, m.PostID, m.Author, m.Content, m.CreatedAt), nil
}

func (r *commentRepository) FindByPostID(ctx context.Context, filter repository.CommentFilter) ([]*entity.Comment, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&commentModel{}).Where("post_id = ?", filter.PostID).Count(&total)

	var models []commentModel
	offset := (filter.Page - 1) * filter.Limit
	if err := r.db.WithContext(ctx).
		Where("post_id = ?", filter.PostID).
		Order("created_at DESC").
		Offset(offset).Limit(filter.Limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	comments := make([]*entity.Comment, len(models))
	for i := range models {
		comments[i] = entity.ReconstituteComment(models[i].ID, models[i].PostID, models[i].Author, models[i].Content, models[i].CreatedAt)
	}
	return comments, total, nil
}

func (r *commentRepository) Save(ctx context.Context, comment *entity.Comment) error {
	return r.db.WithContext(ctx).Create(&commentModel{
		ID:        comment.ID(),
		PostID:    comment.PostID(),
		Author:    comment.Author(),
		Content:   comment.Content(),
		CreatedAt: comment.CreatedAt(),
	}).Error
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&commentModel{}, "id = ?", id).Error
}
