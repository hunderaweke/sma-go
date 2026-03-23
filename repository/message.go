package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	db.AutoMigrate(&domain.Message{})
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(in domain.Message) (*domain.Message, error) {
	if in.RoomId == uuid.Nil {
		return nil, domain.RequiredField("room_id")
	}
	in.FromUnique = strings.TrimSpace(in.FromUnique)
	in.Text = strings.TrimSpace(in.Text)

	if in.FromUnique == "" {
		return nil, domain.RequiredField("from_unique")
	}
	if in.Text == "" {
		return nil, domain.RequiredField("text")
	}

	var room domain.Room
	if err := r.db.First(&room, "id = ?", in.RoomId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("room")
		}
		return nil, err
	}

	var sender domain.Identity
	if err := r.db.Where("unique_string = ?", in.FromUnique).First(&sender).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("identity")
		}
		return nil, err
	}

	if err := r.db.Create(&in).Error; err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *messageRepository) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.RequiredField("id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.InvalidField("id", "must be a valid uuid")
	}

	res := r.db.Delete(&domain.Message{}, "id = ?", uid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.EntityNotFound("message")
	}
	return nil
}

func (r *messageRepository) GetByID(id string) (*domain.Message, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.RequiredField("id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.InvalidField("id", "must be a valid uuid")
	}

	var m domain.Message
	if err := r.db.First(&m, "id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("message")
		}
		return nil, err
	}
	return &m, nil
}

func (r *messageRepository) GetAll(opts options.MessageFetchOptions) ([]domain.Message, error) {
	var items []domain.Message

	base := r.db.Model(&domain.Message{})
	if rcv := strings.TrimSpace(opts.RoomUniqueString); rcv != "" {
		base = base.Joins("JOIN rooms ON rooms.id = messages.room_id").Where("rooms.unique_string = ?", rcv)
	}
	if s := strings.TrimSpace(opts.SenderUniqueString); s != "" {
		base = base.Where("from_unique = ?", s)
	}

	q := base.Order("messages.created_at DESC").Limit(opts.Limit()).Offset(opts.Offset())
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
