package wyvern

type OrderSide int
type SaleKind int
type FeeMethod int
type HowToCall int

const (
	Buy  OrderSide = 0
	Sell OrderSide = 1

	FixedPrice   SaleKind = 0
	DutchAuction SaleKind = 1

	ProtocolFee FeeMethod = 0
	SplitFee    FeeMethod = 1

	Call         HowToCall = 0
	DelegateCall HowToCall = 1
	StaticCall   HowToCall = 2
	Create       HowToCall = 3

	NullAddress = "0x0000000000000000000000000000000000000000"
)

type WyvernNFTAsset struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Quantity string `json:"quantity"`
}

type ExchangeMetadata struct {
	Asset  *WyvernNFTAsset `json:"asset"`
	Schema string          `json:"schema"`
}

type User struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
}

type Account struct {
	ID         int64   `json:"id"`
	Address    string  `json:"address"`
	ImageUrl   *string `json:"profile_img_url"`
	User       *User   `json:"user"`
	IsVerified bool    `json:"is_verified"`
}
