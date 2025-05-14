package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"context"
)

const (
	LIMIT  = 20
	OFFSET = 0
)

func (r *repository) ListTweetsByUserID(ctx context.Context, userID uint64) ([]domain.Tweet, error) {
	var tweets []domain.Tweet
	query := `
		SELECT id, content, user_id, created_at
		FROM tweets
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	limit := LIMIT
	offset := OFFSET
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tweet domain.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.Content, &tweet.UserID, &tweet.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}
