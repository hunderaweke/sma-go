package main

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hunderaweke/sma-go/config"
	_ "github.com/hunderaweke/sma-go/config"
	"github.com/hunderaweke/sma-go/database"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/repository"
	"github.com/hunderaweke/sma-go/server/router"
	"github.com/hunderaweke/sma-go/usecases"
	"github.com/hunderaweke/sma-go/utils"
)

func main() {
	// Ensure .env exists for local setup convenience
	db, err := database.NewDB(database.SQLite)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Identity{}, &domain.Message{}, &domain.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	identityRepo := repository.NewIdentityRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	identityUC := usecases.NewIdentityUsecase(identityRepo)
	pgpHandler := utils.NewPGPHandler()
	messageUC := usecases.NewMessageUsecase(messageRepo, identityUC, pgpHandler)
	userUC := usecases.NewUserUsecase(userRepo)
	analyticsUC := usecases.NewAnalyticsUsecase(analyticsRepo)
	app := router.NewRouter(identityUC, messageUC, analyticsUC, userUC)
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.WebUrl,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
