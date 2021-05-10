package datastore

type DatastoreSort struct {
	Field string
	IsAsc bool
}

type AssetsFilter struct {
	CreatedByID *int64
	Sort        *DatastoreSort
}

type AccountsFilter struct {
	Query *string
	Sort  *DatastoreSort
}
