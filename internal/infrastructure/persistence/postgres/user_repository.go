package postgres

import (
	"context"
	"fmt"
	"meobeo-talk-api/internal/domain/entity"
	"meobeo-talk-api/internal/domain/repository"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
        SELECT id, username, email, display_name, avatar_url, is_active, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	var users []*entity.User
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
        SELECT id, username, email, display_name, avatar_url, is_active, created_at, updated_at
        FROM users
        WHERE id = $1
    `

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
        SELECT id, username, email, password_hash, display_name, avatar_url, is_active, created_at, updated_at
        FROM users
        WHERE email = $1
    `

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
        INSERT INTO users (username, email, password_hash, display_name, avatar_url, is_active)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.PasswordHash, user.DisplayName, user.AvatarURL, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
        UPDATE users
        SET username = $1, email = $2, display_name = $3, avatar_url = $4, is_active = $5, updated_at = NOW()
        WHERE id = $6
        RETURNING updated_at
    `

	err := r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.DisplayName, user.AvatarURL, user.IsActive, user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
