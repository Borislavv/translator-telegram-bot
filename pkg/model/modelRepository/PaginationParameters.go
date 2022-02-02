package modelRepository

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

type PaginationParameters struct {
	page  int
	limit int
}

type WithoutPaginationParameters struct {
}

// NewPaginationParameters - constructor of PaginationParameters struct
func NewPaginationParameters(page int, limit int) *PaginationParameters {
	if page == 0 {
		page = DefaultPage
	}

	if limit == 0 {
		limit = DefaultLimit
	}

	return &PaginationParameters{
		page:  page,
		limit: limit,
	}
}

func NewWithoutPaginationParameters() *WithoutPaginationParameters {
	return &WithoutPaginationParameters{}
}

// GetPage - getter of page prop.
func (pagination *PaginationParameters) GetPage() int {
	return pagination.page
}

// GetLimit - getter of limit prop.
func (pagination *PaginationParameters) GetLimit() int {
	return pagination.limit
}

// NeedPaginate - tell u, need or not paginate by this struct
func (pagination *PaginationParameters) NeedPaginate() bool {
	return true
}

// GetPage - return 0, meaning without pagination
func (withoutPagination *WithoutPaginationParameters) GetPage() int {
	return 0
}

// GetLimit - return 0, meaning without pagination
func (withoutPagination *WithoutPaginationParameters) GetLimit() int {
	return 0
}

// NeedPaginate - tell u, need or not paginate by this struct
func (pagination *WithoutPaginationParameters) NeedPaginate() bool {
	return false
}
