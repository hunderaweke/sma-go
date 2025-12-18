package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerIdentityRoutes(r *gin.Engine, ctrl *controller.IdentityController) {
	identityRoutes := r.Group("/identity")
	identityRoutes.POST("", ctrl.Create)
	identityRoutes.GET("", ctrl.List)
	identityRoutes.GET("/:unique", ctrl.GetByUnique)
	identityRoutes.DELETE("/:id", ctrl.Delete)
}
