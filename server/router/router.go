package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/utils"
)

func NewRouter(identityUC domain.IdentityUsecase, messageUC domain.MessageUsecase, analyticsUC domain.AnalyticsUsecase) *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	identityCtrl := controller.NewIdentityController(identityUC, utils.NewPGPHandler())
	messageCtrl := controller.NewMessageController(messageUC)
	analyticsCtrl := controller.NewAnalyticsController(analyticsUC)

	registerIdentityRoutes(app, identityCtrl)
	registerMessageRoutes(app, messageCtrl)
	registerAnalyticsRoutes(app, analyticsCtrl)

	return app
}
