package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerMessageRoutes(r *gin.Engine, ctrl *controller.MessageController) {
	messageRoutes := r.Group("/messages")
	messageRoutes.POST("", ctrl.Create)
	messageRoutes.GET("/:id", ctrl.GetByID)
	messageRoutes.GET("/receiver/:unique", ctrl.GetByReceiver)
}
