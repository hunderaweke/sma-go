package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/server/controller"
	"github.com/hunderaweke/sma-go/utils"
)

func NewRouter(identityUC domain.IdentityUsecase, messageUC domain.MessageUsecase, analyticsUC domain.AnalyticsUsecase) *gin.Engine {
	r := gin.Default()

	identityCtrl := controller.NewIdentityController(identityUC, utils.NewPGPHandler())
	messageCtrl := controller.NewMessageController(messageUC)
	analyticsCtrl := controller.NewAnalyticsController(analyticsUC)

	registerIdentityRoutes(r, identityCtrl)
	registerMessageRoutes(r, messageCtrl)
	registerAnalyticsRoutes(r, analyticsCtrl)

	return r
}
