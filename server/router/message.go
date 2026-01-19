package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerMessageRoutes(r *fiber.App, ctrl *controller.MessageController) {
	messageRoutes := r.Group("/messages")
	messageRoutes.Post("", ctrl.Create)
	messageRoutes.Get("/:id", ctrl.GetByID)
	messageRoutes.Get("/receiver/:unique", ctrl.GetByReceiver)
	messageRoutes.Get("/receive/:uniqueString", ctrl.ReceiveMessages)
}
