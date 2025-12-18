package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/server/controller"
)

func NewRouter(identityUC domain.IdentityUsecase, messageUC domain.MessageUsecase) *gin.Engine {
	r := gin.Default()

	identityCtrl := controller.NewIdentityController(identityUC)
	messageCtrl := controller.NewMessageController(messageUC)

	registerIdentityRoutes(r, identityCtrl)
	registerMessageRoutes(r, messageCtrl)

	return r
}
