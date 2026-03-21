package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/server/middlewares"
)

func registerRoomRoutes(r *fiber.App, roomCtrl *controller.RoomController, messageCtrl *controller.MessageController) {
	roomRoutes := r.Group("/rooms")
	roomRoutes.Post("", middlewares.JWTMiddleware, roomCtrl.Create)
	roomRoutes.Get("", middlewares.JWTMiddleware, roomCtrl.ListMine)
	roomRoutes.Get("/:uniqueString", middlewares.JWTMiddleware, roomCtrl.GetByUniqueString)
	roomRoutes.Delete("/:uniqueString", middlewares.JWTMiddleware, roomCtrl.Delete)

	messageRoutes := roomRoutes.Group("/:uniqueString/messages")
	messageRoutes.Post("", messageCtrl.CreateInRoom)
	messageRoutes.Get("", middlewares.JWTMiddleware, messageCtrl.ListInRoom)
	messageRoutes.Get("/:id", middlewares.JWTMiddleware, messageCtrl.GetByIDInRoom)
	messageRoutes.Delete("/:id", middlewares.JWTMiddleware, messageCtrl.DeleteInRoom)
}
