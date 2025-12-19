package usecases

import (
	"errors"
	"strings"

	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"github.com/hunderaweke/sma-go/utils"
)

type messageUsecase struct {
	iu         domain.IdentityUsecase
	pgpHandler *utils.PGPHandler
	repo       domain.MessageRepository
}

func NewMessageUsecase(repo domain.MessageRepository, identityUC domain.IdentityUsecase, pgpHandler *utils.PGPHandler) domain.MessageUsecase {
	return &messageUsecase{repo: repo, iu: identityUC, pgpHandler: pgpHandler}
}

func (u *messageUsecase) Create(m domain.Message) (*domain.Message, error) {
	identity, err := u.iu.GetByUniqueString(m.ToUnique)
	if err != nil {
		return nil, err
	}
	if identity.IsPublic == true {
		return u.repo.Create(m)
	}
	publicKey, err := u.pgpHandler.ParsePublicKey(identity.PublicKey)
	if err != nil {
		return nil, err
	}
	encryptedMsg, err := u.pgpHandler.Encrypt(m.Text, publicKey)
	if err != nil {
		return nil, err
	}
	m.Text = encryptedMsg
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
