package service

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/subtrack/backend/internal/model"
	"github.com/subtrack/backend/internal/repository"
	"github.com/subtrack/backend/pkg/password"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrPasswordTooShort   = errors.New("password must be at least 12 characters")
	ErrPasswordTooLong    = errors.New("password must be at most 128 characters")
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, input *model.RegisterInput) (*model.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if !isValidEmail(email) {
		return nil, "", ErrInvalidEmail
	}

	if len(input.Password) < 12 {
		return nil, "", ErrPasswordTooShort
	}
	if len(input.Password) > 128 {
		return nil, "", ErrPasswordTooLong
	}

	hash, err := password.Hash(input.Password)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userRepo.Create(ctx, email, hash)
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	_, token, err := s.sessionRepo.Create(ctx, user.ID, "", "")
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, input *model.LoginInput, ipAddress, userAgent string) (*model.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	valid, err := password.Verify(input.Password, user.PasswordHash)
	if err != nil || !valid {
		return nil, "", ErrInvalidCredentials
	}

	_, token, err := s.sessionRepo.Create(ctx, user.ID, ipAddress, userAgent)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil
	}
	return s.sessionRepo.Delete(ctx, session.ID)
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*model.User, *model.Session, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, nil, err
	}

	s.sessionRepo.ExtendExpiry(ctx, session.ID)

	return user, session, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) UpdatePreferences(ctx context.Context, userID uuid.UUID, prefs *model.UserPreferences) error {
	return s.userRepo.UpdatePreferences(ctx, userID, prefs)
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	valid, err := password.Verify(oldPassword, user.PasswordHash)
	if err != nil || !valid {
		return ErrInvalidCredentials
	}

	if len(newPassword) < 12 {
		return ErrPasswordTooShort
	}
	if len(newPassword) > 128 {
		return ErrPasswordTooLong
	}

	hash, err := password.Hash(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}

	return s.sessionRepo.DeleteAllForUser(ctx, userID)
}

func (s *AuthService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

func isValidEmail(email string) bool {
	if len(email) < 5 || len(email) > 255 {
		return false
	}

	atIndex := strings.Index(email, "@")
	if atIndex < 1 || atIndex > len(email)-3 {
		return false
	}

	domain := email[atIndex+1:]
	if !strings.Contains(domain, ".") {
		return false
	}

	for _, c := range email {
		if unicode.IsSpace(c) {
			return false
		}
	}

	return true
}
