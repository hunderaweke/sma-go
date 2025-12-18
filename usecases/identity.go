package usecases

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type identityUsecase struct {
	repo domain.IdentityRepository
}

func NewIdentityUsecase(repo domain.IdentityRepository) domain.IdentityUsecase {
	return &identityUsecase{repo: repo}
}

func (u *identityUsecase) Create(in domain.Identity) (*domain.Identity, error) {
	in.UniqueString = strings.TrimSpace(in.UniqueString)
	in.PublicKey = strings.TrimSpace(in.PublicKey)
	return u.repo.Create(in)
}

func (u *identityUsecase) Delete(id uuid.UUID) error {
	return u.repo.Delete(id)
}

func (u *identityUsecase) GetByUniqueString(uniqueString string) (*domain.Identity, error) {
	return u.repo.GetByUniqueString(strings.TrimSpace(uniqueString))
}

func (u *identityUsecase) GetAll(opts options.BaseFetchOptions) (domain.MultipleIdentity, error) {
	return u.repo.GetAll(opts)
}
