package controller

import (
	"github.com/gofiber/fiber/v2"
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

func (ic *IdentityController) Create(c *fiber.Ctx) error {
	var req createIdentity
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if ic.pgpHandler == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "encryption not configured"})
	}
	identity, err := ic.usecase.Create(domain.Identity{PublicKey: req.PublicKey, UniqueString: req.UniqueString, IsPublic: req.IsPublic})
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(identity)
}

func (ic *IdentityController) GetByUnique(c *fiber.Ctx) error {
	unique := c.Params("unique")
	identity, err := ic.usecase.GetByUniqueString(unique)
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(identity)
}
func (ic *IdentityController) GetAllIdentities(c *fiber.Ctx) error {
	identities, err := ic.usecase.GetAll(options.BaseFetchOptions{})
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(identities)
}
