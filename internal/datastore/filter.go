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
	Symbol  *string
	Address *string
	Sort    *SortOption
}

type OrderFilter struct {
	Side                 *int
	SaleKind             *int
	PaymentTokenAddress  *string
	AssetContractAddress *string
	TokenID              *int64
	MakerID              *int64
	TakerID              *int64
	Sort                 *SortOption
}
