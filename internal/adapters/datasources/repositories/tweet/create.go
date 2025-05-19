package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"context"
)

func (r *repository) Create(ctx context.Context, tweet *domain.Tweet) (uint64, error) {
	// Prepare the SQL statement with RETURNING id
	query := "INSERT INTO tweets (user_id, content) VALUES ($1, $2) RETURNING id"
	var id uint64
	err := r.db.QueryRowContext(ctx, query, tweet.UserID, tweet.Content).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
