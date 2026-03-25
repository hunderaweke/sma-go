package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hunderaweke/sma-go/config"
	_ "github.com/hunderaweke/sma-go/config"
	"github.com/hunderaweke/sma-go/database"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/repository"
	"github.com/hunderaweke/sma-go/server/router"
	"github.com/hunderaweke/sma-go/usecases"
)

func main() {
	// Ensure .env exists for local setup convenience
	db, err := database.NewDB(database.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Identity{}, &domain.Message{}, &domain.User{}, &domain.Room{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	identityRepo := repository.NewIdentityRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	identityUC := usecases.NewIdentityUsecase(identityRepo)
	messageUC := usecases.NewMessageUsecase(messageRepo)
	userUC := usecases.NewUserUsecase(userRepo)
	roomUC := usecases.NewRoomUsecase(roomRepo)
	analyticsUC := usecases.NewAnalyticsUsecase(analyticsRepo)
	app := router.NewRouter(identityUC, messageUC, analyticsUC, userUC, roomUC)
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.WebUrl,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	address := fmt.Sprintf(":%s", config.ServerPort)
	log.Printf("Server is running on %s", address)
	if err := app.Listen(address); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
