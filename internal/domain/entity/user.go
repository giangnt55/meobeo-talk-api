package entity

import (
	"time"
)

type User struct {
	ID           string  `db:"id" json:"id"`
	Username     string  `db:"username" json:"username"`
	Email        string  `db:"email" json:"email"`
	PasswordHash string  `db:"password_hash" json:"-"`
	DisplayName  *string `db:"display_name" json:"display_name,omitempty"`
	AvatarURL    *string `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio          *string `db:"bio" json:"bio,omitempty"`

	// Status
	IsActive      bool `db:"is_active" json:"is_active"`
	EmailVerified bool `db:"email_verified" json:"email_verified"`

	// Counters
	PostCount      int `db:"post_count" json:"post_count"`
	FollowerCount  int `db:"follower_count" json:"follower_count"`
	FollowingCount int `db:"following_count" json:"following_count"`

	// Timestamps
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	LastSeenAt *time.Time `db:"last_seen_at" json:"last_seen_at,omitempty"`
}
