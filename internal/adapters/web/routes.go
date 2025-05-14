package web

import (
	"github.com/gin-gonic/gin"

	"challenge_be/internal/adapters/web/handlers/follow"
	"challenge_be/internal/adapters/web/handlers/tweet"
	"challenge_be/internal/usecases"
)

func RegisterRoutes(router *gin.Engine, usecases *usecases.UseCases) {
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		tweetGroup := api.Group("/tweets")
		{
			tweetGroup.POST("/", tweet.NewCreateTweetHandler(usecases.CreateTweetUseCase))
			tweetGroup.GET("/timeline/:user_id", tweet.NewGetTimelineHandler(usecases.GetTimelineUseCase))

			// Agregá tus rutas acá, por ejemplo:
			// api.GET("/users", handlers.GetUsers)
		}
		followGroup := api.Group("/follows")
		{
			followGroup.POST("/", follow.NewCreateFollowHandler(usecases.CreateFollowUseCase))

		}
	}
}
