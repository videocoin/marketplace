package datastore

type SortOption struct {
	Field string
	IsAsc bool
}

type AssetsFilter struct {
	Status      *string
	CreatedByID *int64
	CANotNull   *bool
	Sort        *SortOption
}

type AccountsFilter struct {
	Query *string
	Sort  *SortOption
}

type TokensFilter struct {
	Symbol *string
	Sort   *SortOption
}
