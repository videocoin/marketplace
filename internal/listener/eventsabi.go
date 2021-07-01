package listener

const (
	exchangeEventsABI = `
	[
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "hash",
					"type": "bytes32"
				},
				{
					"indexed": false,
					"name": "exchange",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "maker",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "taker",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "makerRelayerFee",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "takerRelayerFee",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "makerProtocolFee",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "takerProtocolFee",
					"type": "uint256"
				},
				{
					"indexed": true,
					"name": "feeRecipient",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "feeMethod",
					"type": "uint8"
				},
				{
					"indexed": false,
					"name": "side",
					"type": "uint8"
				},
				{
					"indexed": false,
					"name": "saleKind",
					"type": "uint8"
				},
				{
					"indexed": false,
					"name": "target",
					"type": "address"
				}
			],
			"name": "OrderApprovedPartOne",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "hash",
					"type": "bytes32"
				},
				{
					"indexed": false,
					"name": "howToCall",
					"type": "uint8"
				},
				{
					"indexed": false,
					"name": "calldata",
					"type": "bytes"
				},
				{
					"indexed": false,
					"name": "replacementPattern",
					"type": "bytes"
				},
				{
					"indexed": false,
					"name": "staticTarget",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "staticExtradata",
					"type": "bytes"
				},
				{
					"indexed": false,
					"name": "paymentToken",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "basePrice",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "extra",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "listingTime",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "expirationTime",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "salt",
					"type": "uint256"
				},
				{
					"indexed": false,
					"name": "orderbookInclusionDesired",
					"type": "bool"
				}
			],
			"name": "OrderApprovedPartTwo",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": true,
					"name": "hash",
					"type": "bytes32"
				}
			],
			"name": "OrderCancelled",
			"type": "event"
		},
		{
			"anonymous": false,
			"inputs": [
				{
					"indexed": false,
					"name": "buyHash",
					"type": "bytes32"
				},
				{
					"indexed": false,
					"name": "sellHash",
					"type": "bytes32"
				},
				{
					"indexed": true,
					"name": "maker",
					"type": "address"
				},
				{
					"indexed": true,
					"name": "taker",
					"type": "address"
				},
				{
					"indexed": false,
					"name": "price",
					"type": "uint256"
				},
				{
					"indexed": true,
					"name": "metadata",
					"type": "bytes32"
				}
			],
			"name": "OrdersMatched",
			"type": "event"
		}
	]
	`
)
