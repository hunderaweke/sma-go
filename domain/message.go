package domain

import (
	"github.com/hunderaweke/sma-go/options"
)

type Message struct {
	Model
	FromUnique string   `gorm:"index;not null"`
	ToUnique   string   `gorm:"index;not null"`
	From       Identity `gorm:"foreignKey:FromUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	To         Identity `gorm:"foreignKey:ToUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Text       string   `gorm:"type:text;not null"`
}

type MultipleMessage struct {
	Meta Pagination
	Data []Message
}

type MessageRepository interface {
	Create(Message) (*Message, error)
	Delete(id string) error
	GetByID(id string) error
	GetAll(opts options.MessageFetchOptions) (MultipleMessage, error)
}

type MessageUsecase interface {
	Create(Message) (*Message, error)
	Delete(id string)
	GetByID(id string)
	GetAll(opts options.MessageFetchOptions) (MultipleMessage, Error)
	GetByReceiverIdentity(recieverID string) (MultipleMessage, Error)
}
