package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	db.AutoMigrate(&domain.User{})
	return &userRepository{db: db}
}

func (r *userRepository) Create(in domain.User) (*domain.User, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Provider = strings.TrimSpace(in.Provider)
	in.ProviderUserID = strings.TrimSpace(in.ProviderUserID)
	in.Email = strings.TrimSpace(in.Email)

	if in.Name == "" {
		return nil, domain.RequiredField("name")
	}
	if in.Provider == "" {
		return nil, domain.RequiredField("provider")
	}
	if in.ProviderUserID == "" {
		return nil, domain.RequiredField("provider_user_id")
	}
	if in.Email == "" {
		return nil, domain.RequiredField("email")
	}

	var exists int64
	if err := r.db.Model(&domain.User{}).
		Where("provider = ? AND provider_user_id = ?", in.Provider, in.ProviderUserID).
		Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, domain.UniqueConstraint("user", "provider_user_id")
	}

	if err := r.db.Create(&in).Error; err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *userRepository) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.RequiredField("id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.InvalidField("id", "must be a valid uuid")
	}

	res := r.db.Delete(&domain.User{}, "id = ?", uid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.EntityNotFound("user")
	}
	return nil
}

func (r *userRepository) Update(id string, data domain.User) (*domain.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.RequiredField("id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.InvalidField("id", "must be a valid uuid")
	}

	data.Name = strings.TrimSpace(data.Name)
	data.Provider = strings.TrimSpace(data.Provider)
	data.ProviderUserID = strings.TrimSpace(data.ProviderUserID)
	data.Email = strings.TrimSpace(data.Email)
	if data.Name == "" {
		return nil, domain.RequiredField("name")
	}
	if data.Provider == "" {
		return nil, domain.RequiredField("provider")
	}
	if data.ProviderUserID == "" {
		return nil, domain.RequiredField("provider_user_id")
	}
	if data.Email == "" {
		return nil, domain.RequiredField("email")
	}

	var exists int64
	if err := r.db.Model(&domain.User{}).
		Where("provider = ? AND provider_user_id = ? AND id <> ?", data.Provider, data.ProviderUserID, uid).
		Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, domain.UniqueConstraint("user", "provider_user_id")
	}

	updates := map[string]any{
		"name":             data.Name,
		"provider":         data.Provider,
		"provider_user_id": data.ProviderUserID,
		"email":            data.Email,
	}

	res := r.db.Model(&domain.User{}).Where("id = ?", uid).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domain.EntityNotFound("user")
	}

	var user domain.User
	if err := r.db.First(&user, "id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("user")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetById(id string) (*domain.User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.RequiredField("id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.InvalidField("id", "must be a valid uuid")
	}

	var user domain.User
	if err := r.db.First(&user, "id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("user")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, domain.RequiredField("email")
	}
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("user")
		}
		return nil, err
	}
	return &user, nil
}
