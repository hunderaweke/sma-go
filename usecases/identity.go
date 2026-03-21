package usecases

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
)

type identityUsecase struct {
	repo domain.IdentityRepository
}

func NewIdentityUsecase(repo domain.IdentityRepository) domain.IdentityUsecase {
	return &identityUsecase{repo: repo}
}

func (u *identityUsecase) Create() (*domain.Identity, error) {
	return u.repo.Create()
}

func (u *identityUsecase) Delete(id uuid.UUID) error {
	return u.repo.Delete(id)
}

func (u *identityUsecase) GetByUniqueString(uniqueString string) (*domain.Identity, error) {
	return u.repo.GetByUniqueString(strings.TrimSpace(uniqueString))
}
