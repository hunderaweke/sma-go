package usecases

import (
	"strings"

	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type roomUsecase struct {
	repo domain.RoomRepository
}

func NewRoomUsecase(repo domain.RoomRepository) domain.RoomUsecase {
	return &roomUsecase{repo: repo}
}

func (u *roomUsecase) Create(in domain.Room) (*domain.Room, error) {
	return u.repo.Create(in)
}

func (u *roomUsecase) Delete(id string) error {
	return u.repo.Delete(strings.TrimSpace(id))
}

func (u *roomUsecase) GetByID(id string) (*domain.Room, error) {
	return u.repo.GetByID(strings.TrimSpace(id))
}

func (u *roomUsecase) GetByUniqueString(uniqueString string) (*domain.Room, error) {
	return u.repo.GetByUniqueString(strings.TrimSpace(uniqueString))
}

func (u *roomUsecase) GetByOwnerId(ownerId string, opts options.BaseFetchOptions) (domain.MultipleRoom, error) {
	return u.repo.GetByOwnerId(strings.TrimSpace(ownerId), opts)
}
