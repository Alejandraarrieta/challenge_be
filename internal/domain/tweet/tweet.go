package domain

import (
	"time"
)

type Tweet struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
