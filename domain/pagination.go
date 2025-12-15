package domain

type Pagination struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	SortBy     string `json:"sort_by,omitempty"`
	SortDesc   bool   `json:"sort_desc,omitempty"`
}
