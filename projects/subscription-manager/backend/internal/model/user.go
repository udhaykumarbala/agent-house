package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID        `json:"id"`
	Email         string           `json:"email"`
	PasswordHash  string           `json:"-"`
	EmailVerified bool             `json:"emailVerified"`
	Preferences   *UserPreferences `json:"preferences"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

type UserPreferences struct {
	Currency           string `json:"currency"`
	Theme              string `json:"theme"`
	ReminderDays       int    `json:"reminderDays"`
	EmailNotifications bool   `json:"emailNotifications"`
}

func DefaultPreferences() *UserPreferences {
	return &UserPreferences{
		Currency:           "USD",
		Theme:              "system",
		ReminderDays:       3,
		EmailNotifications: true,
	}
}

func (p *UserPreferences) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func PreferencesFromJSON(data []byte) (*UserPreferences, error) {
	var p UserPreferences
	if err := json.Unmarshal(data, &p); err != nil {
		return DefaultPreferences(), nil
	}
	return &p, nil
}

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID            uuid.UUID        `json:"id"`
	Email         string           `json:"email"`
	EmailVerified bool             `json:"emailVerified"`
	Preferences   *UserPreferences `json:"preferences"`
	CreatedAt     time.Time        `json:"createdAt"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		Preferences:   u.Preferences,
		CreatedAt:     u.CreatedAt,
	}
}
