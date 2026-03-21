package usecases

import (
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
)

type messageUsecase struct {
	repo domain.MessageRepository
}

func NewMessageUsecase(repo domain.MessageRepository) domain.MessageUsecase {
	return &messageUsecase{repo: repo}
}

func (u *messageUsecase) Create(m domain.Message) (*domain.Message, error) {
	return u.repo.Create(m)
}

func (u *messageUsecase) Delete(id string) error {
	return u.repo.Delete(id)
}

func (u *messageUsecase) GetByID(id string) (*domain.Message, error) {
	return u.repo.GetByID(id)
}

func (u *messageUsecase) GetAll(opts options.MessageFetchOptions) (domain.MultipleMessage, error) {
	return u.repo.GetAll(opts)
}
