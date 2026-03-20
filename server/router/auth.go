package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/server/middlewares"
)

func registerUserRoutes(r *fiber.App, ctrl *controller.UserController) {
	userRoutes := r.Group("/auth")
	userRoutes.Get("/me", middlewares.JWTMiddleware, ctrl.GetMe)
	userRoutes.Get("/:provider", ctrl.SignUpOrLogIn)
	userRoutes.Get("/:provider/callback", ctrl.AuthCallback)
}
