package usecases

import (
	"errors"
	"strings"

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

func (u *messageUsecase) GetAll(opts options.MessageFetchOptions) (domain.MultipleMessage, domain.Error) {
	res, err := u.repo.GetAll(opts)
	if err != nil {
		return domain.MultipleMessage{}, convertError(err, "failed to get messages")
	}
	return res, domain.Error{}
}

func (u *messageUsecase) GetByReceiverIdentity(receiverID string) (domain.MultipleMessage, domain.Error) {
	opts := options.MessageFetchOptions{RoomUniqueString: strings.TrimSpace(receiverID)}
	res, err := u.repo.GetAll(opts)
	if err != nil {
		return domain.MultipleMessage{}, convertError(err, "failed to get messages by receiver")
	}
	return res, domain.Error{}
}

func convertError(err error, msg string) domain.Error {
	if err == nil {
		return domain.Error{}
	}

	var derr *domain.Error
	if errors.As(err, &derr) && derr != nil {
		return *derr
	}
	return *domain.InternalError(err, msg)
}
