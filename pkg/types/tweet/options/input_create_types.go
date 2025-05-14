package tweet

type (
	InputCreateTweet struct {
		Content string `json:"content" validate:"required"`
		UserID  uint64 `json:"user_id" validate:"required"`
	}
)
