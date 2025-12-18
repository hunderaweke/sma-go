package main

import (
	"log"

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

	identityUC := usecases.NewIdentityUsecase(identityRepo)
	pgpHandler := utils.NewPGPHandler()
	messageUC := usecases.NewMessageUsecase(messageRepo, identityUC, pgpHandler)

	r := router.NewRouter(identityUC, messageUC)

	if err := r.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
