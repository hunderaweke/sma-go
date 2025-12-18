package controller

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type MessageController struct {
	usecase domain.MessageUsecase
}

func NewMessageController(uc domain.MessageUsecase) *MessageController {
	return &MessageController{usecase: uc}
}

func (mc *MessageController) Create(c *gin.Context) {
	var req struct {
		FromUnique string `json:"from_unique"`
		ToUnique   string `json:"to_unique"`
		Text       string `json:"text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg, err := mc.usecase.Create(domain.Message{FromUnique: req.FromUnique, ToUnique: req.ToUnique, Text: req.Text})
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(nethttp.StatusCreated, msg)
}

func (mc *MessageController) List(c *gin.Context) {
	var opts options.MessageFetchOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, derr := mc.usecase.GetAll(opts)
	if domainErr := domainErrorFromValue(derr); domainErr != nil {
		writeDomainError(c, domainErr)
		return
	}
	c.JSON(nethttp.StatusOK, res)
}

func (mc *MessageController) GetByReceiver(c *gin.Context) {
	receiver := c.Param("unique")
	res, derr := mc.usecase.GetByReceiverIdentity(receiver)
	if domainErr := domainErrorFromValue(derr); domainErr != nil {
		writeDomainError(c, domainErr)
		return
	}
	c.JSON(nethttp.StatusOK, res)
}

func (mc *MessageController) GetByID(c *gin.Context) {
	id := c.Param("id")
	msg, err := mc.usecase.GetByID(id)
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(nethttp.StatusOK, msg)
}

func (mc *MessageController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := mc.usecase.Delete(id); err != nil {
		writeDomainError(c, err)
		return
	}
	c.Status(nethttp.StatusNoContent)
}
