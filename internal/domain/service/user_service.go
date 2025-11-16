package service

import (
	"context"
	"meobeo-talk-api/internal/dto/request"
	"meobeo-talk-api/internal/dto/response"
)

type UserService interface {
	GetUsers(ctx context.Context, req *request.PaginationRequest) (*response.PaginatedUsers, error)
	GetUserByID(ctx context.Context, id string) (*response.UserResponse, error)
	CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error)
}
