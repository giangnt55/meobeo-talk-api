package entity

import (
	"time"
)

type Post struct {
	ID       string `db:"id" json:"id"`
	AuthorID string `db:"author_id" json:"author_id"`

	// Content
	Title          *string `db:"title" json:"title,omitempty"`
	Content        string  `db:"content" json:"content"`
	ContentPreview *string `db:"content_preview" json:"content_preview,omitempty"`

	// Mood/Emotion
	Mood             *string `db:"mood" json:"mood,omitempty"`
	EmotionIntensity *int    `db:"emotion_intensity" json:"emotion_intensity,omitempty"`

	// Visibility & Status
	Visibility string `db:"visibility" json:"visibility"` // public, followers, private, anonymous
	Status     string `db:"status" json:"status"`         // draft, published, archived

	// Features
	AllowComments bool `db:"allow_comments" json:"allow_comments"`
	IsSensitive   bool `db:"is_sensitive" json:"is_sensitive"`

	// Counters
	CommentCount  int `db:"comment_count" json:"comment_count"`
	ReactionCount int `db:"reaction_count" json:"reaction_count"`
	ViewCount     int `db:"view_count" json:"view_count"`

	// Timestamps
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type PostWithAuthor struct {
	Post
	Author User `db:"author"`
}
