package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/server/middlewares"
)

func registerRoomRoutes(r *fiber.App, ctrl *controller.RoomController) {
	roomRoutes := r.Group("/rooms", middlewares.JWTMiddleware)
	roomRoutes.Post("", ctrl.Create)
	roomRoutes.Get("", ctrl.ListMine)
	roomRoutes.Get("/:uniqueString", ctrl.GetByUniqueString)
	roomRoutes.Delete("/:uniqueString", ctrl.Delete)
}
