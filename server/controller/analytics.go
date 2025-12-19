package controller

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/domain"
)

type AnalyticsController struct {
	usecase domain.AnalyticsUsecase
}

func NewAnalyticsController(uc domain.AnalyticsUsecase) *AnalyticsController {
	return &AnalyticsController{usecase: uc}
}

func (ac *AnalyticsController) Get(c *gin.Context) {
	res, err := ac.usecase.Get()
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(nethttp.StatusOK, res)
}
