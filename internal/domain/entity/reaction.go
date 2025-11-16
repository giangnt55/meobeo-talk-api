package entity

import (
	"time"
)

type Reaction struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	TargetType string    `db:"target_type" json:"target_type"` // post, comment
	TargetID   string    `db:"target_id" json:"target_id"`
	Reaction   string    `db:"reaction" json:"reaction"` // support, hug, love, understand, relate, like
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
