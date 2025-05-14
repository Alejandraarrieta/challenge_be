package follow

type (
	InputCreateFollow struct {
		FollowerID uint64 `json:"follower_id" validate:"required"`
		FolloweeID uint64 `json:"followee_id" validate:"required"`
	}
)
