package entity

import (
	"errors"
	"strings"
	"time"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

var (
	ErrInvalidTitle     = errors.New("title cannot be empty")
	ErrInvalidContent   = errors.New("content cannot be empty")
	ErrInvalidAuthor    = errors.New("author cannot be empty")
	ErrAlreadyPublished = errors.New("post is already published")
	ErrAlreadyDraft     = errors.New("post is already in draft status")
)

type Post struct {
	id        string
	title     string
	content   string
	author    string
	status    PostStatus
	createdAt time.Time
	updatedAt time.Time
}

func NewPost(id, title, content, author string) (*Post, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrInvalidTitle
	}
	if strings.TrimSpace(content) == "" {
		return nil, ErrInvalidContent
	}
	if strings.TrimSpace(author) == "" {
		return nil, ErrInvalidAuthor
	}
	now := time.Now()
	return &Post{
		id:        id,
		title:     strings.TrimSpace(title),
		content:   strings.TrimSpace(content),
		author:    strings.TrimSpace(author),
		status:    PostStatusDraft,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstitute rebuilds a Post from stored data without re-validating.
func Reconstitute(id, title, content, author string, status PostStatus, createdAt, updatedAt time.Time) *Post {
	return &Post{
		id:        id,
		title:     title,
		content:   content,
		author:    author,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (p *Post) ID() string           { return p.id }
func (p *Post) Title() string        { return p.title }
func (p *Post) Content() string      { return p.content }
func (p *Post) Author() string       { return p.author }
func (p *Post) Status() PostStatus   { return p.status }
func (p *Post) CreatedAt() time.Time { return p.createdAt }
func (p *Post) UpdatedAt() time.Time { return p.updatedAt }

func (p *Post) Update(title, content string) error {
	if strings.TrimSpace(title) == "" {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(content) == "" {
		return ErrInvalidContent
	}
	p.title = strings.TrimSpace(title)
	p.content = strings.TrimSpace(content)
	p.updatedAt = time.Now()
	return nil
}

func (p *Post) Publish() error {
	if p.status == PostStatusPublished {
		return ErrAlreadyPublished
	}
	p.status = PostStatusPublished
	p.updatedAt = time.Now()
	return nil
}

func (p *Post) Unpublish() error {
	if p.status == PostStatusDraft {
		return ErrAlreadyDraft
	}
	p.status = PostStatusDraft
	p.updatedAt = time.Now()
	return nil
}
