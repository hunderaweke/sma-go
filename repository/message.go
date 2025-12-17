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
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(in domain.Message) (*domain.Message, error) {
	in.FromUnique = strings.TrimSpace(in.FromUnique)
	in.ToUnique = strings.TrimSpace(in.ToUnique)
	in.Text = strings.TrimSpace(in.Text)

	if in.FromUnique == "" {
		return nil, domain.RequiredField("from_unique")
	}
	if in.ToUnique == "" {
		return nil, domain.RequiredField("to_unique")
	}
	if in.Text == "" {
		return nil, domain.RequiredField("text")
	}

	var sender domain.Identity
	if err := r.db.Where("unique_string = ?", in.FromUnique).First(&sender).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("identity")
		}
		return nil, err
	}

	var recipient domain.Identity
	if err := r.db.Where("unique_string = ?", in.ToUnique).First(&recipient).Error; err != nil {
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

func (r *messageRepository) GetByID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.RequiredField("id")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.InvalidField("id", "must be a valid uuid")
	}

	var m domain.Message
	if err := r.db.First(&m, "id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.EntityNotFound("message")
		}
		return err
	}
	return nil
}

func (r *messageRepository) GetAll(opts options.MessageFetchOptions) (domain.MultipleMessage, error) {
	var (
		items []domain.Message
		total int64
	)

	base := r.db.Model(&domain.Message{})
	if s := strings.TrimSpace(opts.SenderUniqueString); s != "" {
		base = base.Where("from_unique = ?", s)
	}
	if rcv := strings.TrimSpace(opts.RoomUniqueString); rcv != "" {
		base = base.Where("to_unique = ?", rcv)
	}

	if err := base.Count(&total).Error; err != nil {
		return domain.MultipleMessage{}, err
	}

	sortField := sanitizeMessageSortField(opts.SortBy)
	sortDir := "ASC"
	if opts.SortDesc {
		sortDir = "DESC"
	}
	order := sortField + " " + sortDir

	q := base.Order(order).Limit(opts.Limit()).Offset(opts.Offset())
	if err := q.Find(&items).Error; err != nil {
		return domain.MultipleMessage{}, err
	}

	page := opts.GetPage()
	size := opts.GetPageSize()
	totalPages := 0
	if size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}
	meta := domain.Pagination{
		Page:       page,
		PageSize:   size,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
		SortBy:     sortField,
		SortDesc:   opts.SortDesc,
	}

	return domain.MultipleMessage{Meta: meta, Data: items}, nil
}

func sanitizeMessageSortField(in string) string {
	key := strings.TrimSpace(strings.ToLower(in))
	switch key {
	case "id":
		return "id"
	case "from_unique":
		return "from_unique"
	case "to_unique":
		return "to_unique"
	case "updated_at":
		return "updated_at"
	case "created_at", "":
		fallthrough
	default:
		return "created_at"
	}
}
