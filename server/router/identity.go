package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerIdentityRoutes(r *fiber.App, ctrl *controller.IdentityController) {
	identityRoutes := r.Group("/identities")
	identityRoutes.Post("", ctrl.Create)
	identityRoutes.Get("", ctrl.GetAllIdentities)
	identityRoutes.Get("/:unique", ctrl.GetByUnique)
}
