package models

type Identity struct {
	Model
	PublicKey    string `json:"public_key,omitempty" gorm:"not null"`
	UniqueString string `json:"unique_string,omitempty" gorm:"not null;uniqueIndex"`
}

type IdentityRepository interface {
	Create()
	Update()
	Delete()
	GetByID()
	GetAll()
}
