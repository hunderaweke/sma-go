package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/server/controller"
)

func registerIdentityRoutes(r *gin.Engine, ctrl *controller.IdentityController) {
	identityRoutes := r.Group("/identities")
	identityRoutes.POST("", ctrl.Create)
	identityRoutes.GET("/:unique", ctrl.GetByUnique)
}
