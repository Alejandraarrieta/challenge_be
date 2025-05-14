package tweet

import (
	domain "challenge_be/internal/domain/tweet"
	"challenge_be/internal/usecases/tweet"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTimeline godoc
// @Summary Obtener timeline de un usuario
// @Description Devuelve los tweets del timeline del usuario especificado por ID
// @Tags tweets
// @Produce json
// @Param user_id path int true "ID del usuario"
// @Success 200 {object} map[string][]domain.Tweet
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tweets/{user_id} [get]
func NewGetTimelineHandler(usecase tweet.GetTimelineUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDStr := c.Param("user_id")
		if userIDStr == "" {
			c.JSON(400, gin.H{"error": "user_id is required"})
			return
		}
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "user_id must be a number"})
			return
		}
		tweets, err := usecase.Execute(c.Request.Context(), uint64(userID))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to get timeline"})
			return
		}
		if len(tweets) == 0 {
			c.JSON(200, gin.H{"tweets": []domain.Tweet{}})
			return
		}
		c.JSON(200, gin.H{"tweets": tweets})

	}
}
