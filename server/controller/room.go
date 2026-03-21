package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type RoomController struct {
	roomUsecase domain.RoomUsecase
	userUsecase domain.UserUsecase
}

func NewRoomController(roomUC domain.RoomUsecase, userUC domain.UserUsecase) *RoomController {
	return &RoomController{roomUsecase: roomUC, userUsecase: userUC}
}

func (rc *RoomController) Create(c *fiber.Ctx) error {
	owner, err := rc.currentUser(c)
	if err != nil {
		return writeDomainError(c, err)
	}

	room, err := rc.roomUsecase.Create(domain.Room{
		OwnerID: owner.ID,
	})
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(room)
}

func (rc *RoomController) ListMine(c *fiber.Ctx) error {
	owner, err := rc.currentUser(c)
	if err != nil {
		return writeDomainError(c, err)
	}

	var opts options.BaseFetchOptions
	if err := c.QueryParser(&opts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rooms, err := rc.roomUsecase.GetByOwnerId(owner.ID.String(), opts)
	if err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(rooms)
}

func (rc *RoomController) GetByUniqueString(c *fiber.Ctx) error {
	room, err := rc.roomUsecase.GetByUniqueString(c.Params("uniqueString"))
	if err != nil {
		return writeDomainError(c, err)
	}
	if err := rc.ensureOwner(c, room); err != nil {
		return writeDomainError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(room)
}

func (rc *RoomController) Delete(c *fiber.Ctx) error {
	room, err := rc.roomUsecase.GetByUniqueString(c.Params("uniqueString"))
	if err != nil {
		return writeDomainError(c, err)
	}
	if err := rc.ensureOwner(c, room); err != nil {
		return writeDomainError(c, err)
	}
	if err := rc.roomUsecase.Delete(c.Params("uniqueString")); err != nil {
		return writeDomainError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (rc *RoomController) currentUser(c *fiber.Ctx) (*domain.User, error) {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return nil, domain.New(domain.Unauthorized, "missing authenticated user")
	}
	return rc.userUsecase.GetById(userID)
}

func (rc *RoomController) ensureOwner(c *fiber.Ctx, room *domain.Room) error {
	owner, err := rc.currentUser(c)
	if err != nil {
		return err
	}
	if room == nil || room.OwnerID != owner.ID {
		return domain.New(domain.Forbidden, "you do not own this room")
	}
	return nil
}
