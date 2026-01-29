package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subtrack/backend/internal/model"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, userID uuid.UUID, input *model.SubscriptionInput) (*model.Subscription, error) {
	sub := &model.Subscription{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		Currency:     "USD",
		BillingCycle: input.BillingCycle,
		CustomDays:   input.CustomDays,
		Category:     model.CategoryOther,
		Icon:         input.Icon,
		Color:        input.Color,
		Notes:        input.Notes,
		IsActive:     true,
		ReminderDays: 3,
	}

	if input.Currency != nil {
		sub.Currency = *input.Currency
	}
	if input.Category != nil {
		sub.Category = *input.Category
	}
	if input.ReminderDays != nil {
		sub.ReminderDays = *input.ReminderDays
	}

	nextBilling, err := time.Parse("2006-01-02", input.NextBillingDate)
	if err != nil {
		return nil, err
	}
	sub.NextBillingDate = nextBilling

	err = r.pool.QueryRow(ctx, `
		INSERT INTO subscriptions (
			id, user_id, name, description, price, currency, billing_cycle,
			custom_days, next_billing_date, category, icon, color, notes,
			is_active, reminder_days
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING created_at, updated_at
	`, sub.ID, sub.UserID, sub.Name, sub.Description, sub.Price, sub.Currency,
		sub.BillingCycle, sub.CustomDays, sub.NextBillingDate, sub.Category,
		sub.Icon, sub.Color, sub.Notes, sub.IsActive, sub.ReminderDays,
	).Scan(&sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Subscription, error) {
	sub := &model.Subscription{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, name, description, price, currency, billing_cycle,
			custom_days, next_billing_date, category, icon, color, notes,
			is_active, reminder_days, created_at, updated_at
		FROM subscriptions
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&sub.ID, &sub.UserID, &sub.Name, &sub.Description, &sub.Price, &sub.Currency,
		&sub.BillingCycle, &sub.CustomDays, &sub.NextBillingDate, &sub.Category,
		&sub.Icon, &sub.Color, &sub.Notes, &sub.IsActive, &sub.ReminderDays,
		&sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return sub, nil
}

func (r *SubscriptionRepository) GetAllForUser(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, description, price, currency, billing_cycle,
			custom_days, next_billing_date, category, icon, color, notes,
			is_active, reminder_days, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY next_billing_date ASC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.Name, &sub.Description, &sub.Price, &sub.Currency,
			&sub.BillingCycle, &sub.CustomDays, &sub.NextBillingDate, &sub.Category,
			&sub.Icon, &sub.Color, &sub.Notes, &sub.IsActive, &sub.ReminderDays,
			&sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetActiveForUser(ctx context.Context, userID uuid.UUID) ([]*model.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, description, price, currency, billing_cycle,
			custom_days, next_billing_date, category, icon, color, notes,
			is_active, reminder_days, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1 AND is_active = TRUE
		ORDER BY next_billing_date ASC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.Name, &sub.Description, &sub.Price, &sub.Currency,
			&sub.BillingCycle, &sub.CustomDays, &sub.NextBillingDate, &sub.Category,
			&sub.Icon, &sub.Color, &sub.Notes, &sub.IsActive, &sub.ReminderDays,
			&sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetUpcoming(ctx context.Context, userID uuid.UUID, days int) ([]*model.Subscription, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, name, description, price, currency, billing_cycle,
			custom_days, next_billing_date, category, icon, color, notes,
			is_active, reminder_days, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
			AND is_active = TRUE
			AND next_billing_date <= CURRENT_DATE + INTERVAL '1 day' * $2
			AND next_billing_date >= CURRENT_DATE
		ORDER BY next_billing_date ASC
	`, userID, days)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.Name, &sub.Description, &sub.Price, &sub.Currency,
			&sub.BillingCycle, &sub.CustomDays, &sub.NextBillingDate, &sub.Category,
			&sub.Icon, &sub.Color, &sub.Notes, &sub.IsActive, &sub.ReminderDays,
			&sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, id, userID uuid.UUID, input *model.SubscriptionInput) (*model.Subscription, error) {
	nextBilling, err := time.Parse("2006-01-02", input.NextBillingDate)
	if err != nil {
		return nil, err
	}

	currency := "USD"
	if input.Currency != nil {
		currency = *input.Currency
	}

	category := model.CategoryOther
	if input.Category != nil {
		category = *input.Category
	}

	reminderDays := 3
	if input.ReminderDays != nil {
		reminderDays = *input.ReminderDays
	}

	sub := &model.Subscription{}
	err = r.pool.QueryRow(ctx, `
		UPDATE subscriptions SET
			name = $3,
			description = $4,
			price = $5,
			currency = $6,
			billing_cycle = $7,
			custom_days = $8,
			next_billing_date = $9,
			category = $10,
			icon = $11,
			color = $12,
			notes = $13,
			reminder_days = $14
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, name, description, price, currency, billing_cycle,
			custom_days, next_billing_date, category, icon, color, notes,
			is_active, reminder_days, created_at, updated_at
	`, id, userID, input.Name, input.Description, input.Price, currency,
		input.BillingCycle, input.CustomDays, nextBilling, category,
		input.Icon, input.Color, input.Notes, reminderDays,
	).Scan(
		&sub.ID, &sub.UserID, &sub.Name, &sub.Description, &sub.Price, &sub.Currency,
		&sub.BillingCycle, &sub.CustomDays, &sub.NextBillingDate, &sub.Category,
		&sub.Icon, &sub.Color, &sub.Notes, &sub.IsActive, &sub.ReminderDays,
		&sub.CreatedAt, &sub.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return sub, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM subscriptions WHERE id = $1 AND user_id = $2
	`, id, userID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}

func (r *SubscriptionRepository) ToggleActive(ctx context.Context, id, userID uuid.UUID, isActive bool) error {
	result, err := r.pool.Exec(ctx, `
		UPDATE subscriptions SET is_active = $3 WHERE id = $1 AND user_id = $2
	`, id, userID, isActive)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrSubscriptionNotFound
	}

	return nil
}
