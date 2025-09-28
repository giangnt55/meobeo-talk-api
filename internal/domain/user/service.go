package user

import (
	"errors"
	"meobeo-talk-api/internal/dto/request"
	"meobeo-talk-api/internal/dto/response"
	"meobeo-talk-api/internal/models"
	"meobeo-talk-api/internal/utils"

	"github.com/google/uuid"
)

type Service interface {
	CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error)
	UpdateUser(id uuid.UUID, req *request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteUser(id uuid.UUID) error
	GetUserByID(id uuid.UUID) (*response.UserResponse, error)
	GetUsers(pagination *utils.Pagination) (*response.UserListResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error) {
	// Check if email exists
	exists, err := s.repo.Exists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password, // Should be hashed in production
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *service) UpdateUser(id uuid.UUID, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *service) DeleteUser(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("user not found")
	}

	return s.repo.Delete(id)
}

func (s *service) GetUserByID(id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.toUserResponse(user), nil
}

func (s *service) GetUsers(pagination *utils.Pagination) (*response.UserListResponse, error) {
	users, total, err := s.repo.FindAll(pagination)
	if err != nil {
		return nil, err
	}

	return s.toUserListResponse(users, pagination, total), nil
}

// Helper methods to convert models to responses
func (s *service) toUserResponse(u *models.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (s *service) toUserListResponse(users []models.User, pagination *utils.Pagination, total int64) *response.UserListResponse {
	var userResponses []response.UserResponse
	for _, u := range users {
		userResponses = append(userResponses, *s.toUserResponse(&u))
	}

	totalPages := int(total) / pagination.Limit
	if int(total)%pagination.Limit > 0 {
		totalPages++
	}

	return &response.UserListResponse{
		Users: userResponses,
		Meta: response.PaginationMeta{
			CurrentPage: pagination.Page,
			PerPage:     pagination.Limit,
			Total:       total,
			TotalPages:  totalPages,
		},
	}
}
