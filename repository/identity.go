package repository

import (
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"gorm.io/gorm"
)

type identityRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) domain.IdentityRepository {
	db.AutoMigrate(&domain.Identity{})
	return &identityRepository{db: db}
}

func (r *identityRepository) Create() (*domain.Identity, error) {
	newIdent := domain.Identity{}
	if err := r.db.Create(&newIdent).Error; err != nil {
		return nil, err
	}
	newIdent.UniqueString = base64.URLEncoding.EncodeToString(newIdent.ID[:])[:12]
	if err := r.db.Save(&newIdent).Error; err != nil {
		return nil, err
	}
	return &newIdent, nil
}

func (r *identityRepository) Delete(id uuid.UUID) error {
	res := r.db.Delete(&domain.Identity{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.EntityNotFound("identity")
	}
	return nil
}

func (r *identityRepository) GetByUniqueString(uniqueString string) (*domain.Identity, error) {
	uniqueString = strings.TrimSpace(uniqueString)
	if uniqueString == "" {
		return nil, domain.RequiredField("unique_string")
	}
	var ident domain.Identity
	if err := r.db.Where("unique_string = ?", uniqueString).First(&ident).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("identity")
		}
		return nil, err
	}
	return &ident, nil
}
