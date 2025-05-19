package follow

import (
	domain "challenge_be/internal/domain/follow"
	"context"
	errr "database/sql"
)

func (r *repository) Create(ctx context.Context, follow *domain.Follow) error {
	sql := `
		INSERT INTO follows (follower_id, followee_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
		RETURNING id;
	`
	var id int64
	err := r.db.QueryRowContext(ctx, sql, follow.FollowerID, follow.FolloweeID).Scan(&id)
	if err != nil {
		if err == errr.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
