package domain

type Analytics struct {
	Identities int `json:"identities"`
	Messages   int `json:"messages"`
}

type AnalyticsRepository interface {
	Get() (*Analytics, error)
}

type AnalyticsUsecase interface {
	Get() (*Analytics, error)
}
