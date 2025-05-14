package main

// @title Challenge API
// @version 1.0
// @description API para crear, listar y eliminar tweets
// @contact.name Alejandra
// @host localhost:8080
// @BasePath /api

import (
	repositories "challenge_be/internal/adapters/datasources/repositories"
	"challenge_be/internal/adapters/web"
	"challenge_be/internal/platform/cache"
	"challenge_be/internal/platform/database"
	"challenge_be/internal/usecases"
	followuc "challenge_be/internal/usecases/follow"
	tweetuc "challenge_be/internal/usecases/tweet"
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "challenge_be/cmd/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	_ = godotenv.Load()
	// Crear contexto base con cancelación
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatalf("failed to start the application: %v", err)
	}

}

func run(ctx context.Context) error {

	//conetion Postgres
	db, err := database.GetSQLClientInstance()
	if err != nil {
		return err
	}
	defer db.Close()

	//conetion Cache-Redis
	cacheClient, err := cache.GetRedisClientInstance()
	if err != nil {
		return err
	}
	defer cacheClient.Close()

	// Crear instancias de los repositorios
	repos := repositories.CreateRepository(db, cacheClient)

	// Crear instancias de los casos de uso
	useCases := &usecases.UseCases{
		CreateTweetUseCase:  tweetuc.NewCreateUseCase(repos.TweetRepository, repos.TweetCacheRepository),
		CreateFollowUseCase: followuc.NewCreateUseCase(repos.FollowRepository, repos.FollowCacheRepository),
		GetTimelineUseCase:  tweetuc.NewGetTimelineUseCase(repos.TweetRepository, repos.TweetCacheRepository, repos.FollowCacheRepository),
		// Aquí irían otros casos de uso
	}

	// Iniciar motor de rutas con Gin
	router := gin.Default()

	// Agregar ruta Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Registrar rutas
	web.RegisterRoutes(router, useCases)

	// Correr servidor en puerto 8080
	if err := router.Run(":8080"); err != nil {
		return err
	}

	return nil
}
