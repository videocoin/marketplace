package datastore

type SortOption struct {
	Field string
	IsAsc bool
}

type AssetsFilter struct {
	Statuses    []string
	Ids         []int64
	CreatedByID *int64
	OwnerID     *int64
	OnSale      *bool
	Sold        *bool
	Minted      *bool
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
	Ids                  []int64
	Side                 *int
	SaleKind             *int
	PaymentTokenAddress  *string
	AssetContractAddress *string
	TokenID              *int64
	MakerID              *int64
	TakerID              *int64
	Sort                 *SortOption
}

type ActivityFilter struct {
	GroupID     *string
	CreatedByID *int64
	Sort        *SortOption
}
