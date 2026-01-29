package model

import (
	"time"

	"github.com/google/uuid"
)

type BillingCycle string

const (
	BillingWeekly  BillingCycle = "weekly"
	BillingMonthly BillingCycle = "monthly"
	BillingYearly  BillingCycle = "yearly"
	BillingCustom  BillingCycle = "custom"
)

type Category string

const (
	CategoryEntertainment Category = "entertainment"
	CategoryProductivity  Category = "productivity"
	CategoryLifestyle     Category = "lifestyle"
	CategoryUtilities     Category = "utilities"
	CategoryNews          Category = "news"
	CategoryEducation     Category = "education"
	CategoryFinance       Category = "finance"
	CategoryHealth        Category = "health"
	CategoryOther         Category = "other"
)

type Subscription struct {
	ID              uuid.UUID    `json:"id"`
	UserID          uuid.UUID    `json:"userId"`
	Name            string       `json:"name"`
	Description     *string      `json:"description,omitempty"`
	Price           float64      `json:"price"`
	Currency        string       `json:"currency"`
	BillingCycle    BillingCycle `json:"billingCycle"`
	CustomDays      *int         `json:"customDays,omitempty"`
	NextBillingDate time.Time    `json:"nextBillingDate"`
	Category        Category     `json:"category"`
	Icon            *string      `json:"icon,omitempty"`
	Color           *string      `json:"color,omitempty"`
	Notes           *string      `json:"notes,omitempty"`
	IsActive        bool         `json:"isActive"`
	ReminderDays    int          `json:"reminderDays"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}

type SubscriptionInput struct {
	Name            string       `json:"name"`
	Description     *string      `json:"description,omitempty"`
	Price           float64      `json:"price"`
	Currency        *string      `json:"currency,omitempty"`
	BillingCycle    BillingCycle `json:"billingCycle"`
	CustomDays      *int         `json:"customDays,omitempty"`
	NextBillingDate string       `json:"nextBillingDate"`
	Category        *Category    `json:"category,omitempty"`
	Icon            *string      `json:"icon,omitempty"`
	Color           *string      `json:"color,omitempty"`
	Notes           *string      `json:"notes,omitempty"`
	ReminderDays    *int         `json:"reminderDays,omitempty"`
}

func (s *Subscription) MonthlyEquivalent() float64 {
	switch s.BillingCycle {
	case BillingWeekly:
		return s.Price * 52 / 12
	case BillingMonthly:
		return s.Price
	case BillingYearly:
		return s.Price / 12
	case BillingCustom:
		if s.CustomDays != nil && *s.CustomDays > 0 {
			return s.Price * 30 / float64(*s.CustomDays)
		}
		return s.Price
	default:
		return s.Price
	}
}

func (s *Subscription) YearlyEquivalent() float64 {
	switch s.BillingCycle {
	case BillingWeekly:
		return s.Price * 52
	case BillingMonthly:
		return s.Price * 12
	case BillingYearly:
		return s.Price
	case BillingCustom:
		if s.CustomDays != nil && *s.CustomDays > 0 {
			return s.Price * 365 / float64(*s.CustomDays)
		}
		return s.Price * 12
	default:
		return s.Price * 12
	}
}

func (s *Subscription) DaysUntilBilling() int {
	now := time.Now().Truncate(24 * time.Hour)
	billing := s.NextBillingDate.Truncate(24 * time.Hour)
	return int(billing.Sub(now).Hours() / 24)
}

func (s *Subscription) IsDueSoon(days int) bool {
	return s.DaysUntilBilling() <= days && s.DaysUntilBilling() >= 0
}
