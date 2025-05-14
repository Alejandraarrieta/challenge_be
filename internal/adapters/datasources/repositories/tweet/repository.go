package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"database/sql"
)

type (
	repository struct {
		db *sql.DB
	}
)

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{db}
}
