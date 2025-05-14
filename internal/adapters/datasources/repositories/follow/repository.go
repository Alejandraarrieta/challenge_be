package follow

import (
	domain "challenge_be/internal/domain/follow"
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
