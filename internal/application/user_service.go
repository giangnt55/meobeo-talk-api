package application

import (
	"context"
	"fmt"
	"meobeo-talk-api/internal/domain/entity"
	"meobeo-talk-api/internal/domain/repository"
	"meobeo-talk-api/internal/domain/service"
	"meobeo-talk-api/internal/dto/request"
	"meobeo-talk-api/internal/dto/response"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers(ctx context.Context, req *request.PaginationRequest) (*response.PaginatedUsers, error) {
	req.Normalize()

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users
	users, err := s.userRepo.FindAll(ctx, req.PageSize, req.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	// Convert to response
	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = response.ToUserResponse(user)
	}

	return &response.PaginatedUsers{
		Data:       userResponses,
		Pagination: response.NewPagination(req.Page, req.PageSize, total),
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return response.ToUserResponse(user), nil
}

func (s *userService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  &req.FullName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return response.ToUserResponse(user), nil
}
