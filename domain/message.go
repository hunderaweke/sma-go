package domain

import (
	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/options"
)

type Message struct {
	Model
	RoomId     uuid.UUID `json:"room_id,omitempty" gorm:"not null"`
	Room       Room      `gorm:"foreignKey:RoomId;constraint:OnUpdate:SET NULL,OnDelete:SET NULL" json:"-"`
	FromUnique string    `json:"from_unique,omitempty" gorm:"not null"`
	Text       string    `json:"text,omitempty" gorm:"not null"`
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
	GetByRoomUniqueString(roomUniqueString string, opts options.MessageFetchOptions) (MultipleMessage, error)
}
