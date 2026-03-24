package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identity struct {
	Model
	UniqueString string    `json:"unique_string" gorm:"not null;uniqueIndex"`
	ExpireDate   time.Time `json:"expire_date" gorm:"not null"`
}

func (i *Identity) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	i.ExpireDate = time.Now().Add(72 * time.Hour)
	return nil
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
