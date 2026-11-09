package domain

import "time"

type Incident struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Priority    string    `gorm:"size:20;not null" json:"priority"`
	Status      string    `gorm:"size:20;default:'open'" json:"status"`
	UserID      uint      `gorm:"index" json:"user_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
