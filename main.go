package main

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	_ "github.com/hunderaweke/sma-go/config"
	"github.com/hunderaweke/sma-go/database"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/repository"
	"github.com/hunderaweke/sma-go/server/router"
	"github.com/hunderaweke/sma-go/usecases"
	"github.com/hunderaweke/sma-go/utils"
)

func main() {
	// Ensure .env.sample exists for local setup convenience
	db, err := database.NewPostgresConn()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.Identity{}, &domain.Message{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	identityRepo := repository.NewIdentityRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	identityUC := usecases.NewIdentityUsecase(identityRepo)
	pgpHandler := utils.NewPGPHandler()
	messageUC := usecases.NewMessageUsecase(messageRepo, identityUC, pgpHandler)
	analyticsUC := usecases.NewAnalyticsUsecase(analyticsRepo)

	app := router.NewRouter(identityUC, messageUC, analyticsUC)
	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Connection event 1 - Payload: %+v", ep)
	})
	app.Get("/ws/:id", socketio.New(func(kws *socketio.Websocket) {
		uniqueString := kws.Params("unique")
		kws.SetAttribute("unique_string", uniqueString)
		kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", uniqueString, kws.UUID)), true, socketio.TextMessage)
		kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", uniqueString, kws.UUID)), socketio.TextMessage)
	}))
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
