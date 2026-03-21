package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/utils"
)

func NewRouter(identityUC domain.IdentityUsecase, messageUC domain.MessageUsecase, analyticsUC domain.AnalyticsUsecase, userUC domain.UserUsecase, roomUC domain.RoomUsecase) *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	identityCtrl := controller.NewIdentityController(identityUC, utils.NewPGPHandler())
	analyticsCtrl := controller.NewAnalyticsController(analyticsUC)
	authCtrl := controller.NewAuthController(userUC)
	roomCtrl := controller.NewRoomController(roomUC, userUC)
	messageCtrl := controller.NewMessageController(messageUC, roomUC, userUC)

	registerIdentityRoutes(app, identityCtrl)
	registerAnalyticsRoutes(app, analyticsCtrl)
	registerAuthRoutes(app, authCtrl)
	registerRoomRoutes(app, roomCtrl, messageCtrl)

	return app
}
