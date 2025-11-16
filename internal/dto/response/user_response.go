package response

import (
	"meobeo-talk-api/internal/domain/entity"
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Avatar    *string   `json:"avatar"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginatedUsers struct {
	Data       []*UserResponse `json:"data"`
	Pagination *Pagination     `json:"pagination"`
}

func ToUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  *user.DisplayName,
		Avatar:    user.AvatarURL,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}
}
