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

func (ic *IdentityController) Create(c *gin.Context) {
	var query struct {
		IsPublic bool `form:"is_public"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if ic.pgpHandler == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption not configured"})
		return
	}

	key, err := ic.pgpHandler.GenerateKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate key"})
		return
	}

	publicKey, err := key.GetArmoredPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export public key"})
		return
	}

	fp := key.GetFingerprint()
	if len(fp) < 12 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid key fingerprint"})
		return
	}
	uniqueString := fp[:12]
	identity, err := ic.usecase.Create(domain.Identity{PublicKey: publicKey, UniqueString: uniqueString, IsPublic: query.IsPublic})
	if err != nil {
		writeDomainError(c, err)
		return
	}
	privateKey, err := key.Armor()
	var response struct {
		domain.Identity
		PrivateKey string `json:"private_key"`
	}
	response = struct {
		domain.Identity
		PrivateKey string `json:"private_key"`
	}{
		Identity:   *identity,
		PrivateKey: privateKey,
	}
	c.JSON(http.StatusCreated, response)
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
