package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subtrack/backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailExists = errors.New("email already exists")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (*model.User, error) {
	user := &model.User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		Preferences:   model.DefaultPreferences(),
	}

	prefsJSON, err := user.Preferences.ToJSON()
	if err != nil {
		return nil, err
	}

	err = r.pool.QueryRow(ctx, `
		INSERT INTO users (id, email, password_hash, email_verified, preferences)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`, user.ID, user.Email, user.PasswordHash, user.EmailVerified, prefsJSON).Scan(
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrEmailExists
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	var prefsJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, email_verified, preferences, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified,
		&prefsJSON, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Preferences, _ = model.PreferencesFromJSON(prefsJSON)
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	var prefsJSON []byte

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, email_verified, preferences, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.EmailVerified,
		&prefsJSON, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Preferences, _ = model.PreferencesFromJSON(prefsJSON)
	return user, nil
}

func (r *UserRepository) UpdatePreferences(ctx context.Context, id uuid.UUID, prefs *model.UserPreferences) error {
	prefsJSON, err := prefs.ToJSON()
	if err != nil {
		return err
	}

	result, err := r.pool.Exec(ctx, `
		UPDATE users SET preferences = $1 WHERE id = $2
	`, prefsJSON, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE users SET password_hash = $1 WHERE id = $2
	`, passwordHash, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && err.Error() != "" &&
		(contains(err.Error(), "duplicate key") || contains(err.Error(), "UNIQUE constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
