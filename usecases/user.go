package usecases

import (
	"strings"

	"github.com/hunderaweke/sma-go/domain"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) Create(in domain.User) (*domain.User, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Provider = strings.TrimSpace(in.Provider)
	in.ProviderUserID = strings.TrimSpace(in.ProviderUserID)
	in.Email = strings.TrimSpace(in.Email)
	return u.repo.Create(in)
}

func (u *userUsecase) Delete(id string) error {
	return u.repo.Delete(strings.TrimSpace(id))
}

func (u *userUsecase) GetByEmail(email string) (*domain.User, error) {
	email = strings.TrimSpace(email)
	return u.repo.GetByEmail(email)
}

func (u *userUsecase) Update(id string, data domain.User) (*domain.User, error) {
	data.Name = strings.TrimSpace(data.Name)
	data.Provider = strings.TrimSpace(data.Provider)
	data.ProviderUserID = strings.TrimSpace(data.ProviderUserID)
	data.Email = strings.TrimSpace(data.Email)
	return u.repo.Update(strings.TrimSpace(id), data)
}

func (u *userUsecase) GetById(id string) (*domain.User, error) {
	return u.repo.GetById(strings.TrimSpace(id))
}
