package datastore

type SortOption struct {
	Field string
	IsAsc bool
}

type AssetsFilter struct {
	CreatedByID *int64
	Sort        *SortOption
}

type AccountsFilter struct {
	Query *string
	Sort  *SortOption
}
