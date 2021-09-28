package model

import (
	"github.com/gocraft/dbr/v2"
	"time"
)

const (
	ActivityGroupBids      = "bids"
	ActivityGroupPurchases = "purchases"
	ActivityGroupSales     = "sales"

	ActivityTypeVIDReceived = "vid_received"
	ActivityTypePurchased   = "purchased"
	ActivityTypeSold        = "sold"
)

type Activity struct {
	ID          int64         `db:"id"`
	IsNew       bool          `db:"is_new"`
	CreatedAt   *time.Time    `db:"created_at"`
	CreatedByID int64         `db:"created_by_id"`
	TypeID      string        `db:"type_id"`
	GroupID     string        `db:"group_id"`
	AssetID     dbr.NullInt64 `db:"asset_id"`
	OrderID     dbr.NullInt64 `db:"order_id"`

	CreatedBy *Account `db:"-"`
	Asset     *Asset   `db:"-"`
	Order     *Order   `db:"-"`
}
