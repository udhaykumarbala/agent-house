package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/subtrack/backend/internal/model"
	"github.com/subtrack/backend/internal/repository"
)

var (
	ErrInvalidName        = errors.New("name must be 1-100 characters")
	ErrInvalidPrice       = errors.New("price must be between 0 and 999999.99")
	ErrInvalidCycle       = errors.New("invalid billing cycle")
	ErrInvalidCategory    = errors.New("invalid category")
	ErrInvalidDate        = errors.New("invalid date format")
	ErrInvalidColor       = errors.New("invalid color format")
	ErrInvalidReminder    = errors.New("reminder days must be 0-30")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

var validCycles = map[model.BillingCycle]bool{
	model.BillingWeekly:  true,
	model.BillingMonthly: true,
	model.BillingYearly:  true,
	model.BillingCustom:  true,
}

var validCategories = map[model.Category]bool{
	model.CategoryEntertainment: true,
	model.CategoryProductivity:  true,
	model.CategoryLifestyle:     true,
	model.CategoryUtilities:     true,
	model.CategoryNews:          true,
	model.CategoryEducation:     true,
	model.CategoryFinance:       true,
	model.CategoryHealth:        true,
	model.CategoryOther:         true,
}

var colorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, userID uuid.UUID, input *model.SubscriptionInput) (*model.Subscription, error) {
	if err := s.validate(input); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, userID, input)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) GetAll(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	return s.repo.GetAllForUser(ctx, userID)
}

func (s *SubscriptionService) GetActive(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	return s.repo.GetActiveForUser(ctx, userID)
}

func (s *SubscriptionService) GetUpcoming(ctx context.Context, userID uuid.UUID, days int) ([]*model.Subscription, error) {
	if days < 1 {
		days = 30
	}
	if days > 365 {
		days = 365
	}
	return s.repo.GetUpcoming(ctx, userID, days)
}

func (s *SubscriptionService) Update(ctx context.Context, id, userID uuid.UUID, input *model.SubscriptionInput) (*model.Subscription, error) {
	if err := s.validate(input); err != nil {
		return nil, err
	}

	sub, err := s.repo.Update(ctx, id, userID, input)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}
	return nil
}

func (s *SubscriptionService) validate(input *model.SubscriptionInput) error {
	name := strings.TrimSpace(input.Name)
	if len(name) < 1 || len(name) > 100 {
		return ErrInvalidName
	}
	input.Name = name

	if input.Price < 0 || input.Price > 999999.99 {
		return ErrInvalidPrice
	}

	if !validCycles[input.BillingCycle] {
		return ErrInvalidCycle
	}

	if input.Category != nil && !validCategories[*input.Category] {
		return ErrInvalidCategory
	}

	_, err := time.Parse("2006-01-02", input.NextBillingDate)
	if err != nil {
		return ErrInvalidDate
	}

	if input.Color != nil && *input.Color != "" && !colorRegex.MatchString(*input.Color) {
		return ErrInvalidColor
	}

	if input.ReminderDays != nil && (*input.ReminderDays < 0 || *input.ReminderDays > 30) {
		return ErrInvalidReminder
	}

	return nil
}
