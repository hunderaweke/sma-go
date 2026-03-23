package domain

import (
	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/options"
)

type Room struct {
	Model
	UniqueString string    `json:"unique_string,omitempty" gorm:"not null;uniqueIndex"`
	Name         string    `json:"name,omitempty" gorm:"not null"`
	OwnerID      uuid.UUID `json:"owner_id,omitempty" gorm:"not null"`
	Owner        User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL" json:"-"`
}

type MultipleRoom struct {
	Meta Pagination `json:"meta,omitempty"`
	Data []Room     `json:"data,omitempty"`
}

type RoomRepository interface {
	Create(Room) (*Room, error)
	Delete(id string) error
	GetByID(id string) (*Room, error)
	UpdateName(id string, name string) (*Room, error)
	GetByUniqueString(uniqueString string) (*Room, error)
	GetByOwnerId(ownerId string, opts options.BaseFetchOptions) (MultipleRoom, error)
}

type RoomUsecase interface {
	Create(Room) (*Room, error)
	Delete(id string) error
	GetByID(id string) (*Room, error)
	GetByUniqueString(uniqueString string) (*Room, error)
	GetByOwnerId(ownerId string, opts options.BaseFetchOptions) (MultipleRoom, error)
	UpdateName(id string, name string) (*Room, error)
}
