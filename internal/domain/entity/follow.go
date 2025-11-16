package entity

import "time"

type Follow struct {
	FollowerID string    `db:"follower_id" json:"follower_id"`
	FolloweeID string    `db:"followee_id" json:"followee_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type FollowWithUser struct {
	Follow
	User User `db:"user"`
}
