package user

import (
	"meobeo-talk-api/internal/models"
	"meobeo-talk-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll(pagination *utils.Pagination) ([]models.User, int64, error)
	Exists(email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *repository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindAll(pagination *utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Get total count
	query.Count(&total)

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Limit(pagination.Limit).Offset(offset).Find(&users).Error

	return users, total, err
}

func (r *repository) Exists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
