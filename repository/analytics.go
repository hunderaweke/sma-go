package repository

import (
	"github.com/hunderaweke/sma-go/domain"
	"gorm.io/gorm"
)

type analyticsRepo struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) domain.AnalyticsRepository {
	return &analyticsRepo{db: db}
}

func (r *analyticsRepo) Get() (*domain.Analytics, error) {
	var a domain.Analytics
	query := `
		SELECT
			(SELECT COUNT(*) FROM identities) AS identities,
			(SELECT COUNT(*) FROM messages)   AS messages;
	`
	if err := r.db.Raw(query).Scan(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
