package domain

type User struct {
	Model
	Name           string `json:"name,omitempty"`
	Provider       string `json:"provider,omitempty"`
	ProviderUserID string `json:"provider_user_id,omitempty"`
	Email          string `json:"email,omitempty"`
}

type UserRepository interface {
	Create(User) (*User, error)
	Delete(id string) error
	Update(id string, data User) (*User, error)
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}
type UserUsecase interface {
	Create(User) (*User, error)
	Delete(id string) error
	Update(id string, data User) (*User, error)
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
}
