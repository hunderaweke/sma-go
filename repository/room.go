package repository

import (
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

type roomWithMessagesCount struct {
	domain.Room
	MessagesCnt int `gorm:"column:messages_cnt"`
}

func NewRoomRepository(db *gorm.DB) domain.RoomRepository {
	db.AutoMigrate(&domain.Room{})
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(in domain.Room) (*domain.Room, error) {
	if in.OwnerID == uuid.Nil {
		return nil, domain.RequiredField("owner_id")
	}

	var exists int64
	if err := r.db.Model(&domain.Room{}).
		Where("unique_string = ?", in.UniqueString).
		Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, domain.UniqueConstraint("room", "unique_string")
	}

	if err := r.db.Create(&in).Error; err != nil {
		return nil, err
	}
	in.UniqueString = base64.RawURLEncoding.EncodeToString(in.ID[:])[:12]
	if in.Name == "" {
		in.Name = in.UniqueString
	}
	if err := r.db.Save(&in).Error; err != nil {
		return nil, err
	}
	return &in, nil
}

func (r *roomRepository) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.RequiredField("unique_string")
	}

	res := r.db.Delete(&domain.Room{}, "unique_string = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.EntityNotFound("room")
	}
	return nil
}

func (r *roomRepository) GetByID(id string) (*domain.Room, error) {
	return r.GetByUniqueString(id)
}

func (r *roomRepository) UpdateName(id string, name string) (*domain.Room, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, domain.RequiredField("unique_string")
	}

	room, err := r.GetByUniqueString(id)
	if err != nil {
		return nil, err
	}

	room.Name = name
	if err := r.db.Save(room).Error; err != nil {
		return nil, err
	}

	return room, nil
}

func (r *roomRepository) GetByUniqueString(uniqueString string) (*domain.Room, error) {
	uniqueString = strings.TrimSpace(uniqueString)
	if uniqueString == "" {
		return nil, domain.RequiredField("unique_string")
	}
	var row roomWithMessagesCount
	if err := roomQuery(r.db).Where("unique_string = ?", uniqueString).First(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.EntityNotFound("room")
		}
		return nil, err
	}
	row.Room.MessagesCnt = row.MessagesCnt
	return &row.Room, nil
}

func (r *roomRepository) GetByOwnerId(ownerId string, opts options.BaseFetchOptions) (domain.MultipleRoom, error) {
	var (
		items []domain.Room
		total int64
	)

	ownerId = strings.TrimSpace(ownerId)
	if ownerId == "" {
		return domain.MultipleRoom{}, domain.RequiredField("owner_id")
	}

	uid, err := uuid.Parse(ownerId)
	if err != nil {
		return domain.MultipleRoom{}, domain.InvalidField("owner_id", "must be a valid uuid")
	}

	if err := r.db.Model(&domain.Room{}).Where("owner_id = ?", uid).Count(&total).Error; err != nil {
		return domain.MultipleRoom{}, err
	}

	sortField := sanitizeRoomSortField(opts.SortBy)
	sortDir := "ASC"
	if opts.SortDesc {
		sortDir = "DESC"
	}
	order := sortField + " " + sortDir

	var rows []roomWithMessagesCount
	q := roomQuery(r.db).Where("owner_id = ?", uid).Order(order).Limit(opts.Limit()).Offset(opts.Offset())
	if err := q.Find(&rows).Error; err != nil {
		return domain.MultipleRoom{}, err
	}

	items = make([]domain.Room, 0, len(rows))
	for _, row := range rows {
		row.Room.MessagesCnt = row.MessagesCnt
		items = append(items, row.Room)
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

	return domain.MultipleRoom{Meta: meta, Data: items}, nil
}

func roomQuery(db *gorm.DB) *gorm.DB {
	return db.Model(&domain.Room{}).
		Select("rooms.*, (SELECT count(id) FROM messages WHERE messages.room_id = rooms.id) as messages_cnt")
}

func sanitizeRoomSortField(in string) string {
	key := strings.TrimSpace(strings.ToLower(in))
	switch key {
	case "id":
		return "id"
	case "name":
		return "name"
	case "unique_string":
		return "unique_string"
	case "owner_id":
		return "owner_id"
	case "updated_at":
		return "updated_at"
	case "created_at", "":
		fallthrough
	default:
		return "created_at"
	}
}
