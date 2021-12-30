package modelInterface

type PaginationInterface interface {
	GetPage() int
	GetLimit() int
	NeedPaginate() bool
}
