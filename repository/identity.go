package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hunderaweke/sma-go/domain"
	"github.com/hunderaweke/sma-go/options"
	"gorm.io/gorm"
)

type identityRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) domain.IdentityRepository {
	db.AutoMigrate(&domain.Identity{})
	return &identityRepository{db: db}
}

func (r *identityRepository) Create(in domain.Identity) (*domain.Identity, error) {
	in.UniqueString = strings.TrimSpace(in.UniqueString)
	in.PublicKey = strings.TrimSpace(in.PublicKey)
	if in.UniqueString == "" {
		return nil, domain.RequiredField("unique_string")
	}
	if in.PublicKey == "" {
		return nil, domain.RequiredField("public_key")
	}

	var exists int64
	if err := r.db.Model(&domain.Identity{}).
		Where("unique_string = ?", in.UniqueString).
		Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, domain.UniqueConstraint("identity", "unique_string")
	}

	if err := r.db.Create(&in).Error; err != nil {
		return nil, err
	}
	return &in, nil
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

func (r *identityRepository) GetAll(opts options.BaseFetchOptions) (domain.MultipleIdentity, error) {
	var (
		items []domain.Identity
		total int64
	)

	if err := r.db.Model(&domain.Identity{}).Count(&total).Error; err != nil {
		return domain.MultipleIdentity{}, err
	}

	sortField := sanitizeSortField(opts.SortBy)
	sortDir := "ASC"
	if opts.SortDesc {
		sortDir = "DESC"
	}
	order := sortField + " " + sortDir

	q := r.db.Model(&domain.Identity{}).Order(order).Limit(opts.Limit()).Offset(opts.Offset())
	if err := q.Find(&items).Error; err != nil {
		return domain.MultipleIdentity{}, err
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

	return domain.MultipleIdentity{Meta: meta, Data: items}, nil
}

func sanitizeSortField(in string) string {
	key := strings.TrimSpace(strings.ToLower(in))
	switch key {
	case "id":
		return "id"
	case "unique_string":
		return "unique_string"
	case "updated_at":
		return "updated_at"
	case "created_at", "":
		fallthrough
	default:
		return "created_at"
	}
}
