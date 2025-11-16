package entity

import "time"

type Comment struct {
	ID       string  `db:"id" json:"id"`
	PostID   string  `db:"post_id" json:"post_id"`
	AuthorID string  `db:"author_id" json:"author_id"`
	ParentID *string `db:"parent_id" json:"parent_id,omitempty"` // For threading

	Content       string `db:"content" json:"content"`
	ReactionCount int    `db:"reaction_count" json:"reaction_count"`
	Status        string `db:"status" json:"status"` // visible, hidden, deleted

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CommentWithAuthor struct {
	Comment
	Author User `db:"author"`
}
