package domain

import (
	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/options"
)

type Identity struct {
	Model
	PublicKey    string `json:"public_key,omitempty" gorm:"not null"`
	UniqueString string `json:"unique_string,omitempty" gorm:"not null;uniqueIndex"`
}

type MultipleIdentity struct {
	Meta Pagination `json:"meta,omitempty"`
	Data []Identity `json:"data,omitempty"`
}
type IdentityRepository interface {
	Create(Identity) (*Identity, error)
	Delete(id uuid.UUID) error
	GetByUniqueString(uniqueString string) error
	GetAll(opts options.BaseFetchOptions) (MultipleIdentity, error)
}
