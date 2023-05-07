package domain

type Paginator struct {
	Limit      uint
	Offset     uint
	TotalItems uint
}

type Response[T any] struct {
	Paginator Paginator
	Items     []T
}

func NewResponse[RepoModel any, DomainModel any](filter ViewFilter, total int64, items []RepoModel,
	repoToDomain func([]RepoModel) []DomainModel) Response[DomainModel] {

	limit, _ := filter.Limit()
	offset, _ := filter.Offset()
	return Response[DomainModel]{
		Paginator: Paginator{
			Limit:      limit,
			Offset:     offset,
			TotalItems: uint(total),
		},
		Items: repoToDomain(items),
	}
}
