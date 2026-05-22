package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ic-plus-backend/models"
)

// UserRepository handles database operations for the users table.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user and returns the created user.
func (r *UserRepository) Create(ctx context.Context, email, passwordHash, role string) (*models.User, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)`,
		email, passwordHash, role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.FindByID(ctx, int(id))
}

// FindByEmail retrieves a user by email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return user, nil
}

// FindByID retrieves a user by ID.
func (r *UserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return user, nil
}

// StoreOTP stores the OTP code and its expiration for a user.
func (r *UserRepository) StoreOTP(ctx context.Context, userID int, otpCode string, expiredAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET reset_otp_code = ?, reset_otp_expired_at = ? WHERE id = ?`,
		otpCode, expiredAt, userID,
	)
	return err
}

// FindByEmailAndOTP retrieves a user by email and matching non-expired OTP code.
func (r *UserRepository) FindByEmailAndOTP(ctx context.Context, email, otpCode string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users 
		 WHERE email = ? AND reset_otp_code = ? AND reset_otp_expired_at > NOW()`,
		email, otpCode,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email and OTP: %w", err)
	}
	return user, nil
}

// FindByIDAndOTP retrieves a user by ID and matching non-expired OTP code.
func (r *UserRepository) FindByIDAndOTP(ctx context.Context, userID int, otpCode string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users 
		 WHERE id = ? AND reset_otp_code = ? AND reset_otp_expired_at > NOW()`,
		userID, otpCode,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID and OTP: %w", err)
	}
	return user, nil
}

// UpdatePassword updates the user's password hash and clears the OTP.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ?, reset_otp_code = NULL, reset_otp_expired_at = NULL WHERE id = ?`,
		passwordHash, userID,
	)
	return err
}
