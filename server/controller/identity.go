package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"github.com/hunderaweke/sma-go/utils"
)

type IdentityController struct {
	pgpHandler *utils.PGPHandler
	usecase    domain.IdentityUsecase
}

func NewIdentityController(uc domain.IdentityUsecase, handler *utils.PGPHandler) *IdentityController {
	if handler == nil {
		handler = utils.NewPGPHandler()
	}
	return &IdentityController{usecase: uc, pgpHandler: handler}
}

type createIdentity struct {
	PublicKey    string `json:"public_key"`
	IsPublic     bool   `json:"is_public"`
	UniqueString string `json:"unique_string"`
}

func (ic *IdentityController) Create(c *gin.Context) {
	var req createIdentity
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ic.pgpHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption not configured"})
		return
	}
	identity, err := ic.usecase.Create(domain.Identity{PublicKey: req.PublicKey, UniqueString: req.UniqueString, IsPublic: req.IsPublic})
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(http.StatusCreated, identity)
}

func (ic *IdentityController) List(c *gin.Context) {
	var opts options.BaseFetchOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ic.usecase.GetAll(opts)
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (ic *IdentityController) GetByUnique(c *gin.Context) {
	unique := c.Param("unique")
	identity, err := ic.usecase.GetByUniqueString(unique)
	if err != nil {
		writeDomainError(c, err)
		return
	}
	c.JSON(http.StatusOK, identity)
}

func (ic *IdentityController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	if err := ic.usecase.Delete(uid); err != nil {
		writeDomainError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
