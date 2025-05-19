package tweet

type TweetResponse struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
