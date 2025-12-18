package controller

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type IdentityController struct {
	usecase domain.IdentityUsecase
}

func NewIdentityController(uc domain.IdentityUsecase) *IdentityController {
	return &IdentityController{usecase: uc}
}

func (ic *IdentityController) Create(c *gin.Context) {
	var req struct {
		PublicKey    string `json:"public_key"`
		UniqueString string `json:"unique_string"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity, err := ic.usecase.Create(domain.Identity{PublicKey: req.PublicKey, UniqueString: req.UniqueString})
	if err != nil {
		writeDomainError(c, err)
		return
	}

	c.JSON(nethttp.StatusCreated, identity)
}

func (ic *IdentityController) List(c *gin.Context) {
	var opts options.BaseFetchOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ic.usecase.GetAll(opts)
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(nethttp.StatusOK, res)
}

func (ic *IdentityController) GetByUnique(c *gin.Context) {
	unique := c.Param("unique")
	identity, err := ic.usecase.GetByUniqueString(unique)
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(nethttp.StatusOK, identity)
}

func (ic *IdentityController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	if err := ic.usecase.Delete(uid); err != nil {
		writeDomainError(c, err)
		return
	}
	c.Status(nethttp.StatusNoContent)
}
