package listener

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	orderApprovedPartOne = crypto.Keccak256Hash([]byte("OrderApprovedPartOne(bytes32,address,address,address,uint256,uint256,uint256,uint256,address,uint8,uint8,uint8,address)"))
	orderApprovedPartTwo = crypto.Keccak256Hash([]byte("OrderApprovedPartTwo(bytes32,uint8,bytes,bytes,address,bytes,address,uint256,uint256,uint256,uint256,uint256,bool)"))
	orderCanceled        = crypto.Keccak256Hash([]byte("OrderCancelled(bytes32)"))
	ordersMatched        = crypto.Keccak256Hash([]byte("OrdersMatched(bytes32,bytes32,address,address,uint256,bytes32)"))
)

const (
	OrderApproved  = iota
	OrderCancelled = iota
	OrdersMatched  = iota
)

type OrderEvent struct {
	Type     int
	Hash     common.Hash
	SellHash common.Hash
	BuyHash  common.Hash
	Maker    common.Address
	Taker    common.Address
}

type ordersMatchedEvent struct {
	BuyHash  common.Hash
	SellHash common.Hash
	Maker    common.Address
	Taker    common.Address
	Price    *big.Int
	Metadata common.Hash
}

type orderCanceledEvent struct {
	Hash common.Hash
}

type orderApprovedPartOneEvent struct {
	Hash             common.Hash
	Exchange         common.Address
	Maker            common.Address
	Taker            common.Address
	MakerRelayerFee  *big.Int
	TakerRelayerFee  *big.Int
	MakerProtocolFee *big.Int
	TakerProtocolFee *big.Int
	FeeRecipient     common.Address
	FeeMethod        uint8
	Side             uint8
	SaleKind         uint8
	Target           common.Address
}

type orderApprovedPartTwoEvent struct {
	Hash                      common.Hash
	HowToCall                 uint8
	Calldata                  []byte
	ReplacementPattern        []byte
	StaticTarget              common.Address
	StaticExtradata           []byte
	PaymentToken              common.Address
	BasePrice                 *big.Int
	Extra                     *big.Int
	ListingTime               *big.Int
	ExpirationTime            *big.Int
	Salt                      *big.Int
	OrderbookInclusionDesired bool
}
