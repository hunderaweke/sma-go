package options

type BaseFetchOptions struct {
	Page     int    `json:"page" form:"page" query:"page"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" form:"sort_by" query:"sort_by"`
	SortDesc bool   `json:"sort_desc" form:"sort_desc" query:"sort_desc"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func (o BaseFetchOptions) GetPage() int {
	if o.Page < 1 {
		return DefaultPage
	}
	return o.Page
}

func (o BaseFetchOptions) GetPageSize() int {
	size := o.PageSize
	if size <= 0 {
		size = DefaultPageSize
	}
	if size > MaxPageSize {
		size = MaxPageSize
	}
	return size
}

func (o BaseFetchOptions) Offset() int {
	return (o.GetPage() - 1) * o.GetPageSize()
}

func (o BaseFetchOptions) Limit() int {
	return o.GetPageSize()
}

type MessageFetchOptions struct {
	BaseFetchOptions
	RoomUniqueString   string `json:"room_unique_string" form:"room_unique_string" query:"room_unique_string"`
	SenderUniqueString string `json:"sender_unique_string" form:"sender_unique_string" query:"sender_unique_string"`
}
