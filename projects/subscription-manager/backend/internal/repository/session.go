package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subtrack/backend/internal/model"
)

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionExpired = errors.New("session expired")

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Create(ctx context.Context, userID uuid.UUID, ipAddress, userAgent string) (*model.Session, string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, "", err
	}
	token := hex.EncodeToString(tokenBytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	session := &model.Session{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO sessions (id, user_id, token_hash, expires_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`, session.ID, session.UserID, session.TokenHash, session.ExpiresAt,
		session.IPAddress, session.UserAgent).Scan(&session.CreatedAt)

	if err != nil {
		return nil, "", err
	}

	return session, token, nil
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*model.Session, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	session := &model.Session{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, token_hash, expires_at, ip_address, user_agent, created_at
		FROM sessions WHERE token_hash = $1
	`, tokenHash).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt,
		&session.IPAddress, &session.UserAgent, &session.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.IsExpired() {
		r.Delete(ctx, session.ID)
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (r *SessionRepository) ExtendExpiry(ctx context.Context, id uuid.UUID) error {
	newExpiry := time.Now().Add(7 * 24 * time.Hour)
	_, err := r.pool.Exec(ctx, `
		UPDATE sessions SET expires_at = $1 WHERE id = $2
	`, newExpiry, id)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func (r *SessionRepository) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *SessionRepository) GetActiveSessionsForUser(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, token_hash, expires_at, ip_address, user_agent, created_at
		FROM sessions
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.Session
	for rows.Next() {
		session := &model.Session{}
		err := rows.Scan(
			&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt,
			&session.IPAddress, &session.UserAgent, &session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
