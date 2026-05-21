package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidCommentContent   = errors.New("comment content cannot be empty")
	ErrInvalidCommentAuthor    = errors.New("comment author cannot be empty")
	ErrCommentPostNotPublished = errors.New("cannot add comment to an unpublished post")
)

type Comment struct {
	id        string
	postID    string
	author    string
	content   string
	createdAt time.Time
}

func NewComment(id, postID, author, content string) (*Comment, error) {
	if strings.TrimSpace(author) == "" {
		return nil, ErrInvalidCommentAuthor
	}
	if strings.TrimSpace(content) == "" {
		return nil, ErrInvalidCommentContent
	}
	return &Comment{
		id:        id,
		postID:    postID,
		author:    strings.TrimSpace(author),
		content:   strings.TrimSpace(content),
		createdAt: time.Now(),
	}, nil
}

func ReconstituteComment(id, postID, author, content string, createdAt time.Time) *Comment {
	return &Comment{
		id:        id,
		postID:    postID,
		author:    author,
		content:   content,
		createdAt: createdAt,
	}
}

func (c *Comment) ID() string           { return c.id }
func (c *Comment) PostID() string       { return c.postID }
func (c *Comment) Author() string       { return c.author }
func (c *Comment) Content() string      { return c.content }
func (c *Comment) CreatedAt() time.Time { return c.createdAt }
