package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/subtrack/backend/internal/model"
	"github.com/subtrack/backend/internal/repository"
)

type Summary struct {
	TotalMonthly  float64           `json:"totalMonthly"`
	TotalYearly   float64           `json:"totalYearly"`
	ActiveCount   int               `json:"activeCount"`
	ByCategory    map[string]float64 `json:"byCategory"`
	DueSoon       int               `json:"dueSoon"`
	DueSoonAmount float64           `json:"dueSoonAmount"`
}

type UpcomingGroup struct {
	Label         string                 `json:"label"`
	Subscriptions []*model.Subscription `json:"subscriptions"`
	Total         float64               `json:"total"`
}

type AnalyticsService struct {
	subRepo *repository.SubscriptionRepository
}

func NewAnalyticsService(subRepo *repository.SubscriptionRepository) *AnalyticsService {
	return &AnalyticsService{subRepo: subRepo}
}

func (s *AnalyticsService) GetSummary(ctx context.Context, userID uuid.UUID) (*Summary, error) {
	subs, err := s.subRepo.GetActiveForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary := &Summary{
		ByCategory: make(map[string]float64),
	}

	for _, sub := range subs {
		monthly := sub.MonthlyEquivalent()
		yearly := sub.YearlyEquivalent()

		summary.TotalMonthly += monthly
		summary.TotalYearly += yearly
		summary.ActiveCount++

		category := string(sub.Category)
		summary.ByCategory[category] += monthly

		if sub.IsDueSoon(7) {
			summary.DueSoon++
			summary.DueSoonAmount += sub.Price
		}
	}

	summary.TotalMonthly = roundToTwoDecimals(summary.TotalMonthly)
	summary.TotalYearly = roundToTwoDecimals(summary.TotalYearly)
	summary.DueSoonAmount = roundToTwoDecimals(summary.DueSoonAmount)

	for cat := range summary.ByCategory {
		summary.ByCategory[cat] = roundToTwoDecimals(summary.ByCategory[cat])
	}

	return summary, nil
}

func (s *AnalyticsService) GetUpcoming(ctx context.Context, userID uuid.UUID, days int) ([]*UpcomingGroup, error) {
	subs, err := s.subRepo.GetUpcoming(ctx, userID, days)
	if err != nil {
		return nil, err
	}

	thisWeek := &UpcomingGroup{Label: "This Week", Subscriptions: make([]*model.Subscription, 0)}
	nextWeek := &UpcomingGroup{Label: "Next Week", Subscriptions: make([]*model.Subscription, 0)}
	later := &UpcomingGroup{Label: "Later This Month", Subscriptions: make([]*model.Subscription, 0)}

	for _, sub := range subs {
		daysUntil := sub.DaysUntilBilling()
		switch {
		case daysUntil <= 7:
			thisWeek.Subscriptions = append(thisWeek.Subscriptions, sub)
			thisWeek.Total += sub.Price
		case daysUntil <= 14:
			nextWeek.Subscriptions = append(nextWeek.Subscriptions, sub)
			nextWeek.Total += sub.Price
		default:
			later.Subscriptions = append(later.Subscriptions, sub)
			later.Total += sub.Price
		}
	}

	thisWeek.Total = roundToTwoDecimals(thisWeek.Total)
	nextWeek.Total = roundToTwoDecimals(nextWeek.Total)
	later.Total = roundToTwoDecimals(later.Total)

	groups := make([]*UpcomingGroup, 0, 3)
	if len(thisWeek.Subscriptions) > 0 {
		groups = append(groups, thisWeek)
	}
	if len(nextWeek.Subscriptions) > 0 {
		groups = append(groups, nextWeek)
	}
	if len(later.Subscriptions) > 0 {
		groups = append(groups, later)
	}

	return groups, nil
}

func roundToTwoDecimals(val float64) float64 {
	return float64(int(val*100+0.5)) / 100
}
