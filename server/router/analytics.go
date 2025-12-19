package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerAnalyticsRoutes(r *gin.Engine, ctrl *controller.AnalyticsController) {
	analyticsRoutes := r.Group("/analytics")
	analyticsRoutes.GET("", ctrl.Get)
}
