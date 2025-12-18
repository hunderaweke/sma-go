package domain

import (
	"github.com/hunderaweke/sma-go/options"
)

type Message struct {
	Model
	FromUnique string   `gorm:"index;not null" json:"from_unique,omitempty"`
	ToUnique   string   `gorm:"index;not null" json:"to_unique,omitempty"`
	From       Identity `gorm:"foreignKey:FromUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	To         Identity `gorm:"foreignKey:ToUnique;references:UniqueString;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Text       string   `gorm:"type:text;not null" json:"text,omitempty"`
}

type MultipleMessage struct {
	Meta Pagination `json:"meta,omitempty"`
	Data []Message  `json:"data,omitempty"`
}

type MessageRepository interface {
	Create(Message) (*Message, error)
	Delete(id string) error
	GetByID(id string) (*Message, error)
	GetAll(opts options.MessageFetchOptions) (MultipleMessage, error)
}

type MessageUsecase interface {
	Create(Message) (*Message, error)
	Delete(id string) error
	GetByID(id string) (*Message, error)
	GetAll(opts options.MessageFetchOptions) (MultipleMessage, Error)
	GetByReceiverIdentity(recieverID string) (MultipleMessage, Error)
}
