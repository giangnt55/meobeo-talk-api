package entity

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID        string         `db:"id" json:"id"`
	UserID    string         `db:"user_id" json:"user_id"`
	ActorID   *string        `db:"actor_id" json:"actor_id,omitempty"`
	Type      string         `db:"type" json:"type"`       // new_comment, new_reaction, new_follower, mention
	Payload   sql.NullString `db:"payload" json:"payload"` // JSONB stored as string
	IsRead    bool           `db:"is_read" json:"is_read"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
}

type NotificationWithActor struct {
	Notification
	Actor *User `db:"actor"`
}
