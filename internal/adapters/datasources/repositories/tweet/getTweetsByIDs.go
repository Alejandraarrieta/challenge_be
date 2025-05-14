package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"context"
	"fmt"
	"strings"
)

func (r *repository) GetTweetsByIDs(ctx context.Context, ids []uint64) ([]domain.Tweet, error) {
	if len(ids) == 0 {
		return []domain.Tweet{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, content, created_at
		FROM tweets
		WHERE user_id IN (%s)
		ORDER BY created_at DESC
	`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []domain.Tweet
	for rows.Next() {
		var tweet domain.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Content, &tweet.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil

}
