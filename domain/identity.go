package domain

import (
	"github.com/google/uuid"
)

type Identity struct {
	Model
	UniqueString string `json:"unique_string" gorm:"not null;uniqueIndex"`
}

type MultipleIdentity struct {
	Meta Pagination `json:"meta,omitempty"`
	Data []Identity `json:"data,omitempty"`
}
type IdentityRepository interface {
	Create() (*Identity, error)
	Delete(id uuid.UUID) error
	GetByUniqueString(uniqueString string) (*Identity, error)
}
type IdentityUsecase interface {
	Create() (*Identity, error)
	Delete(id uuid.UUID) error
	GetByUniqueString(uniqueString string) (*Identity, error)
}
