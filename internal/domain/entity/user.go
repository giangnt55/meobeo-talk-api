package entity

import "time"

type User struct {
	ID        int64      `db:"id" json:"id"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	Password  string     `db:"password" json:"-"`
	FullName  string     `db:"full_name" json:"full_name"`
	Avatar    *string    `db:"avatar" json:"avatar"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
}
