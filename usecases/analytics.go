package usecases

import "github.com/hunderaweke/sma-go/domain"

type analyticsUsecase struct {
	repo domain.AnalyticsRepository
}

func NewAnalyticsUsecase(repo domain.AnalyticsRepository) domain.AnalyticsUsecase {
	return &analyticsUsecase{repo: repo}
}
func (u *analyticsUsecase) Get() (*domain.Analytics, error) {
	return u.repo.Get()
}
