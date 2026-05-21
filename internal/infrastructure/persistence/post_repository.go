package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"checkapp/internal/domain/entity"
	"checkapp/internal/domain/repository"
)

type postModel struct {
	ID        string `gorm:"primaryKey;type:varchar(36)"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	Author    string `gorm:"not null"`
	Status    string `gorm:"not null;default:'draft'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (postModel) TableName() string { return "posts" }

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	db.AutoMigrate(&postModel{}) //nolint:errcheck
	return &postRepository{db: db}
}

func (r *postRepository) FindByID(ctx context.Context, id string) (*entity.Post, error) {
	var m postModel
	result := r.db.WithContext(ctx).First(&m, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, repository.ErrPostNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return toEntity(&m), nil
}

func (r *postRepository) FindAll(ctx context.Context, filter repository.PostFilter) ([]*entity.Post, int64, error) {
	var total int64
	countQ := r.db.WithContext(ctx).Model(&postModel{})
	if filter.Status != nil {
		countQ = countQ.Where("status = ?", string(*filter.Status))
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []postModel
	findQ := r.db.WithContext(ctx).Model(&postModel{})
	if filter.Status != nil {
		findQ = findQ.Where("status = ?", string(*filter.Status))
	}
	offset := (filter.Page - 1) * filter.Limit
	if err := findQ.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	posts := make([]*entity.Post, len(models))
	for i := range models {
		posts[i] = toEntity(&models[i])
	}
	return posts, total, nil
}

func (r *postRepository) Save(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Create(fromEntity(post)).Error
}

func (r *postRepository) Update(ctx context.Context, post *entity.Post) error {
	return r.db.WithContext(ctx).Save(fromEntity(post)).Error
}

func (r *postRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&postModel{}, "id = ?", id).Error
}

func toEntity(m *postModel) *entity.Post {
	return entity.Reconstitute(m.ID, m.Title, m.Content, m.Author, entity.PostStatus(m.Status), m.CreatedAt, m.UpdatedAt)
}

func fromEntity(p *entity.Post) *postModel {
	return &postModel{
		ID:        p.ID(),
		Title:     p.Title(),
		Content:   p.Content(),
		Author:    p.Author(),
		Status:    string(p.Status()),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
	}
}
