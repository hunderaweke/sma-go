package domain

import (
	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/options"
)

type Message struct {
	Model
	RoomId     uuid.UUID `json:"room_id,omitempty" gorm:"not null"`
	Room       Room      `gorm:"foreignKey:RoomId;constraint:OnUpdate:SET NULL,OnDelete:SET NULL" json:"-"`
	FromUnique string    `gorm:"index;not null" json:"from_unique"`
	Text       string    `gorm:"type:text;not null" json:"text,omitempty"`
}

type MessageRepository interface {
	Create(Message) (*Message, error)
	Delete(id string) error
	GetByID(id string) (*Message, error)
	GetAll(opts options.MessageFetchOptions) ([]Message, error)
}

type MessageUsecase interface {
	Create(Message) (*Message, error)
	Delete(id string) error
	GetByID(id string) (*Message, error)
	GetAll(opts options.MessageFetchOptions) ([]Message, error)
}
