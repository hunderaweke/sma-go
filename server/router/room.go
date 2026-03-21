package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/server/middlewares"
)

func registerRoomRoutes(r *fiber.App, roomCtrl *controller.RoomController, messageCtrl *controller.MessageController) {
	roomRoutes := r.Group("/rooms", middlewares.JWTMiddleware)
	roomRoutes.Post("", roomCtrl.Create)
	roomRoutes.Get("", roomCtrl.ListMine)
	roomRoutes.Get("/:uniqueString", roomCtrl.GetByUniqueString)
	roomRoutes.Delete("/:uniqueString", roomCtrl.Delete)

	messageRoutes := roomRoutes.Group("/:uniqueString/messages")
	messageRoutes.Post("", messageCtrl.CreateInRoom)
	messageRoutes.Get("", messageCtrl.ListInRoom)
	messageRoutes.Get("/receive", messageCtrl.ReceiveMessages)
	messageRoutes.Get("/:id", messageCtrl.GetByIDInRoom)
	messageRoutes.Delete("/:id", messageCtrl.DeleteInRoom)
}
