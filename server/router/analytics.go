package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerAnalyticsRoutes(r *fiber.App, ctrl *controller.AnalyticsController) {
	analyticsRoutes := r.Group("/analytics")
	analyticsRoutes.Get("", ctrl.Get)
}
