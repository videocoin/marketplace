// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nft

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Nft1155ABI is the input ABI used to generate the binding from.
const Nft1155ABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"uri_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Nft1155Bin is the compiled bytecode used for deploying new contracts.
var Nft1155Bin = "0x60806040523480156200001157600080fd5b506040516200268738038062002687833981810160405260208110156200003757600080fd5b81019080805160405193929190846401000000008211156200005857600080fd5b838201915060208201858111156200006f57600080fd5b82518660018202830111640100000000821117156200008d57600080fd5b8083526020830192505050908051906020019080838360005b83811015620000c3578082015181840152602081019050620000a6565b50505050905090810190601f168015620000f15780820380516001836020036101000a031916815260200191505b5060405250505080620001116301ffc9a760e01b6200015a60201b60201c565b62000122816200026360201b60201c565b6200013a63d9b67a2660e01b6200015a60201b60201c565b62000152630e89341c60e01b6200015a60201b60201c565b505062000335565b63ffffffff60e01b817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161415620001f7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f4552433136353a20696e76616c696420696e746572666163652069640000000081525060200191505060405180910390fd5b6001600080837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b80600390805190602001906200027b9291906200027f565b5050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282620002b7576000855562000303565b82601f10620002d257805160ff191683800117855562000303565b8280016001018555821562000303579182015b8281111562000302578251825591602001919060010190620002e5565b5b50905062000312919062000316565b5090565b5b808211156200033157600081600090555060010162000317565b5090565b61234280620003456000396000f3fe608060405234801561001057600080fd5b506004361061009d5760003560e01c806340c10f191161006657806340c10f191461049f5780634e1273f414610503578063a22cb465146106a4578063e985e9c5146106f4578063f242432a1461076e5761009d565b8062fdd58e146100a257806301ffc9a7146101045780630e89341c14610167578063156e29f61461020e5780632eb2c2d61461027c575b600080fd5b6100ee600480360360408110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061087d565b6040518082815260200191505060405180910390f35b61014f6004803603602081101561011a57600080fd5b8101908080357bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916906020019092919050505061095d565b60405180821515815260200191505060405180910390f35b6101936004803603602081101561017d57600080fd5b81019080803590602001909291905050506109c4565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101d35780820151818401526020810190506101b8565b50505050905090810190601f1680156102005780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102646004803603606081101561022457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610a68565b60405180821515815260200191505060405180910390f35b61049d600480360360a081101561029257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156102ef57600080fd5b82018360208201111561030157600080fd5b8035906020019184602083028401116401000000008311171561032357600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561038357600080fd5b82018360208201111561039557600080fd5b803590602001918460208302840111640100000000831117156103b757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561041757600080fd5b82018360208201111561042957600080fd5b8035906020019184600183028401116401000000008311171561044b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610a90565b005b6104eb600480360360408110156104b557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610f1b565b60405180821515815260200191505060405180910390f35b61064d6004803603604081101561051957600080fd5b810190808035906020019064010000000081111561053657600080fd5b82018360208201111561054857600080fd5b8035906020019184602083028401116401000000008311171561056a57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290803590602001906401000000008111156105ca57600080fd5b8201836020820111156105dc57600080fd5b803590602001918460208302840111640100000000831117156105fe57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050610f43565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610690578082015181840152602081019050610675565b505050509050019250505060405180910390f35b6106f2600480360360408110156106ba57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803515159060200190929190505050611055565b005b6107566004803603604081101561070a57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506111ee565b60405180821515815260200191505060405180910390f35b61087b600480360360a081101561078457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190803590602001906401000000008111156107f557600080fd5b82018360208201111561080757600080fd5b8035906020019184600183028401116401000000008311171561082957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050611282565b005b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610904576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b81526020018061219d602b913960400191505060405180910390fd5b6001600083815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000806000837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200190815260200160002060009054906101000a900460ff169050919050565b606060038054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a5c5780601f10610a3157610100808354040283529160200191610a5c565b820191906000526020600020905b815481529060010190602001808311610a3f57829003601f168201915b50505050509050919050565b6000610a85848484604051806020016040528060008152506115f7565b600190509392505050565b8151835114610aea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806122c46028913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415610b70576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806121f16025913960400191505060405180910390fd5b610b786117fa565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161480610bbe5750610bbd85610bb86117fa565b6111ee565b5b610c13576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260328152602001806122166032913960400191505060405180910390fd5b6000610c1d6117fa565b9050610c2d818787878787611802565b60005b8451811015610dfe576000858281518110610c4757fe5b602002602001015190506000858381518110610c5f57fe5b60200260200101519050610ce6816040518060600160405280602a8152602001612248602a91396001600086815260200190815260200160002060008d73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461180a9092919063ffffffff16565b6001600084815260200190815260200160002060008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610d9d816001600085815260200190815260200160002060008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118c490919063ffffffff16565b6001600084815260200190815260200160002060008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505050806001019050610c30565b508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb8787604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b83811015610eae578082015181840152602081019050610e93565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610ef0578082015181840152602081019050610ed5565b5050505090500194505050505060405180910390a4610f1381878787878761194c565b505050505050565b6000610f3983836001604051806020016040528060008152506115f7565b6001905092915050565b60608151835114610f9f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602981526020018061229b6029913960400191505060405180910390fd5b6000835167ffffffffffffffff81118015610fb957600080fd5b50604051908082528060200260200182016040528015610fe85781602001602082028036833780820191505090505b50905060005b845181101561104a5761102785828151811061100657fe5b602002602001015185838151811061101a57fe5b602002602001015161087d565b82828151811061103357fe5b602002602001018181525050806001019050610fee565b508091505092915050565b8173ffffffffffffffffffffffffffffffffffffffff166110746117fa565b73ffffffffffffffffffffffffffffffffffffffff1614156110e1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806122726029913960400191505060405180910390fd5b80600260006110ee6117fa565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff1661119b6117fa565b73ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c318360405180821515815260200191505060405180910390a35050565b6000600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415611308576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001806121f16025913960400191505060405180910390fd5b6113106117fa565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614806113565750611355856113506117fa565b6111ee565b5b6113ab576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806121c86029913960400191505060405180910390fd5b60006113b56117fa565b90506113d58187876113c688611cdb565b6113cf88611cdb565b87611802565b611452836040518060600160405280602a8152602001612248602a91396001600088815260200190815260200160002060008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461180a9092919063ffffffff16565b6001600086815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550611509836001600087815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118c490919063ffffffff16565b6001600086815260200190815260200160002060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f628787604051808381526020018281526020019250505060405180910390a46115ef818787878787611d4c565b505050505050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141561167d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806122ec6021913960400191505060405180910390fd5b60006116876117fa565b90506116a88160008761169988611cdb565b6116a288611cdb565b87611802565b61170b836001600087815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546118c490919063ffffffff16565b6001600086815260200190815260200160002060008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508473ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f628787604051808381526020018281526020019250505060405180910390a46117f381600087878787611d4c565b5050505050565b600033905090565b505050505050565b60008383111582906118b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561187c578082015181840152602081019050611861565b50505050905090810190601f1680156118a95780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5082840390509392505050565b600080828401905083811015611942576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b61196b8473ffffffffffffffffffffffffffffffffffffffff16612059565b15611cd3578373ffffffffffffffffffffffffffffffffffffffff1663bc197c8187878686866040518663ffffffff1660e01b8152600401808673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff168152602001806020018060200180602001848103845287818151815260200191508051906020019060200280838360005b83811015611a23578082015181840152602081019050611a08565b50505050905001848103835286818151815260200191508051906020019060200280838360005b83811015611a65578082015181840152602081019050611a4a565b50505050905001848103825285818151815260200191508051906020019080838360005b83811015611aa4578082015181840152602081019050611a89565b50505050905090810190601f168015611ad15780820380516001836020036101000a031916815260200191505b5098505050505050505050602060405180830381600087803b158015611af657600080fd5b505af1925050508015611b2a57506040513d6020811015611b1657600080fd5b810190808051906020019092919050505060015b611c3457611b3661208a565b80611b415750611be3565b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611ba8578082015181840152602081019050611b8d565b50505050905090810190601f168015611bd55780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260348152602001806121416034913960400191505060405180910390fd5b63bc197c8160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614611cd1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806121756028913960400191505060405180910390fd5b505b505050505050565b60606000600167ffffffffffffffff81118015611cf757600080fd5b50604051908082528060200260200182016040528015611d265781602001602082028036833780820191505090505b5090508281600081518110611d3757fe5b60200260200101818152505080915050919050565b611d6b8473ffffffffffffffffffffffffffffffffffffffff16612059565b15612051578373ffffffffffffffffffffffffffffffffffffffff1663f23a6e6187878686866040518663ffffffff1660e01b8152600401808673ffffffffffffffffffffffffffffffffffffffff1681526020018573ffffffffffffffffffffffffffffffffffffffff16815260200184815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b83811015611e24578082015181840152602081019050611e09565b50505050905090810190601f168015611e515780820380516001836020036101000a031916815260200191505b509650505050505050602060405180830381600087803b158015611e7457600080fd5b505af1925050508015611ea857506040513d6020811015611e9457600080fd5b810190808051906020019092919050505060015b611fb257611eb461208a565b80611ebf5750611f61565b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015611f26578082015181840152602081019050611f0b565b50505050905090810190601f168015611f535780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260348152602001806121416034913960400191505060405180910390fd5b63f23a6e6160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19161461204f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260288152602001806121756028913960400191505060405180910390fd5b505b505050505050565b600080823b905060008111915050919050565b6000601f19601f8301169050919050565b60008160e01c9050919050565b600060443d101561209a5761213d565b60046000803e6120ab60005161207d565b6308c379a081146120bc575061213d565b60405160043d036004823e80513d602482011167ffffffffffffffff821117156120e85750505061213d565b808201805167ffffffffffffffff81111561210757505050505061213d565b8060208301013d85018111156121225750505050505061213d565b61212b8261206c565b60208401016040528296505050505050505b9056fe455243313135353a207472616e7366657220746f206e6f6e2045524331313535526563656976657220696d706c656d656e746572455243313135353a204552433131353552656365697665722072656a656374656420746f6b656e73455243313135353a2062616c616e636520717565727920666f7220746865207a65726f2061646472657373455243313135353a2063616c6c6572206973206e6f74206f776e6572206e6f7220617070726f766564455243313135353a207472616e7366657220746f20746865207a65726f2061646472657373455243313135353a207472616e736665722063616c6c6572206973206e6f74206f776e6572206e6f7220617070726f766564455243313135353a20696e73756666696369656e742062616c616e636520666f72207472616e73666572455243313135353a2073657474696e6720617070726f76616c2073746174757320666f722073656c66455243313135353a206163636f756e747320616e6420696473206c656e677468206d69736d61746368455243313135353a2069647320616e6420616d6f756e7473206c656e677468206d69736d61746368455243313135353a206d696e7420746f20746865207a65726f2061646472657373a2646970667358221220519c1c830ab0fc92729699d368b08a37bb97e1751bae93491d47a5a8b14bb5d164736f6c63430007060033"

// DeployNft1155 deploys a new Ethereum contract, binding an instance of Nft1155 to it.
func DeployNft1155(auth *bind.TransactOpts, backend bind.ContractBackend, uri_ string) (common.Address, *types.Transaction, *Nft1155, error) {
	parsed, err := abi.JSON(strings.NewReader(Nft1155ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Nft1155Bin), backend, uri_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Nft1155{Nft1155Caller: Nft1155Caller{contract: contract}, Nft1155Transactor: Nft1155Transactor{contract: contract}, Nft1155Filterer: Nft1155Filterer{contract: contract}}, nil
}

// Nft1155 is an auto generated Go binding around an Ethereum contract.
type Nft1155 struct {
	Nft1155Caller     // Read-only binding to the contract
	Nft1155Transactor // Write-only binding to the contract
	Nft1155Filterer   // Log filterer for contract events
}

// Nft1155Caller is an auto generated read-only Go binding around an Ethereum contract.
type Nft1155Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Nft1155Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Nft1155Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Nft1155Session struct {
	Contract     *Nft1155          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Nft1155CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Nft1155CallerSession struct {
	Contract *Nft1155Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Nft1155TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Nft1155TransactorSession struct {
	Contract     *Nft1155Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Nft1155Raw is an auto generated low-level Go binding around an Ethereum contract.
type Nft1155Raw struct {
	Contract *Nft1155 // Generic contract binding to access the raw methods on
}

// Nft1155CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Nft1155CallerRaw struct {
	Contract *Nft1155Caller // Generic read-only contract binding to access the raw methods on
}

// Nft1155TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Nft1155TransactorRaw struct {
	Contract *Nft1155Transactor // Generic write-only contract binding to access the raw methods on
}

// NewNft1155 creates a new instance of Nft1155, bound to a specific deployed contract.
func NewNft1155(address common.Address, backend bind.ContractBackend) (*Nft1155, error) {
	contract, err := bindNft1155(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Nft1155{Nft1155Caller: Nft1155Caller{contract: contract}, Nft1155Transactor: Nft1155Transactor{contract: contract}, Nft1155Filterer: Nft1155Filterer{contract: contract}}, nil
}

// NewNft1155Caller creates a new read-only instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Caller(address common.Address, caller bind.ContractCaller) (*Nft1155Caller, error) {
	contract, err := bindNft1155(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Nft1155Caller{contract: contract}, nil
}

// NewNft1155Transactor creates a new write-only instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Transactor(address common.Address, transactor bind.ContractTransactor) (*Nft1155Transactor, error) {
	contract, err := bindNft1155(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Nft1155Transactor{contract: contract}, nil
}

// NewNft1155Filterer creates a new log filterer instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Filterer(address common.Address, filterer bind.ContractFilterer) (*Nft1155Filterer, error) {
	contract, err := bindNft1155(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Nft1155Filterer{contract: contract}, nil
}

// bindNft1155 binds a generic wrapper to an already deployed contract.
func bindNft1155(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Nft1155ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nft1155 *Nft1155Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nft1155.Contract.Nft1155Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nft1155 *Nft1155Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.Contract.Nft1155Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nft1155 *Nft1155Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nft1155.Contract.Nft1155Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nft1155 *Nft1155CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nft1155.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nft1155 *Nft1155TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nft1155 *Nft1155TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nft1155.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Caller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Session) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.BalanceOf(&_Nft1155.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.BalanceOf(&_Nft1155.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155Caller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155Session) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Nft1155.Contract.BalanceOfBatch(&_Nft1155.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155CallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Nft1155.Contract.BalanceOfBatch(&_Nft1155.CallOpts, accounts, ids)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Nft1155 *Nft1155Caller) IsApprovedForAll(opts *bind.CallOpts, account common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "isApprovedForAll", account, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Nft1155 *Nft1155Session) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Nft1155.Contract.IsApprovedForAll(&_Nft1155.CallOpts, account, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address account, address operator) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) IsApprovedForAll(account common.Address, operator common.Address) (bool, error) {
	return _Nft1155.Contract.IsApprovedForAll(&_Nft1155.CallOpts, account, operator)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nft1155.Contract.SupportsInterface(&_Nft1155.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nft1155.Contract.SupportsInterface(&_Nft1155.CallOpts, interfaceId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_Nft1155 *Nft1155Caller) Uri(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "uri", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_Nft1155 *Nft1155Session) Uri(arg0 *big.Int) (string, error) {
	return _Nft1155.Contract.Uri(&_Nft1155.CallOpts, arg0)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 ) view returns(string)
func (_Nft1155 *Nft1155CallerSession) Uri(arg0 *big.Int) (string, error) {
	return _Nft1155.Contract.Uri(&_Nft1155.CallOpts, arg0)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 tokenId, uint256 amount) returns(bool)
func (_Nft1155 *Nft1155Transactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mint", to, tokenId, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 tokenId, uint256 amount) returns(bool)
func (_Nft1155 *Nft1155Session) Mint(to common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint(&_Nft1155.TransactOpts, to, tokenId, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 tokenId, uint256 amount) returns(bool)
func (_Nft1155 *Nft1155TransactorSession) Mint(to common.Address, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint(&_Nft1155.TransactOpts, to, tokenId, amount)
}

// Mint0 is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns(bool)
func (_Nft1155 *Nft1155Transactor) Mint0(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mint0", to, tokenId)
}

// Mint0 is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns(bool)
func (_Nft1155 *Nft1155Session) Mint0(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint0(&_Nft1155.TransactOpts, to, tokenId)
}

// Mint0 is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 tokenId) returns(bool)
func (_Nft1155 *Nft1155TransactorSession) Mint0(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint0(&_Nft1155.TransactOpts, to, tokenId)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_Nft1155 *Nft1155Transactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_Nft1155 *Nft1155Session) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeBatchTransferFrom(&_Nft1155.TransactOpts, from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_Nft1155 *Nft1155TransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeBatchTransferFrom(&_Nft1155.TransactOpts, from, to, ids, amounts, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_Nft1155 *Nft1155Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "safeTransferFrom", from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_Nft1155 *Nft1155Session) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeTransferFrom(&_Nft1155.TransactOpts, from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_Nft1155 *Nft1155TransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeTransferFrom(&_Nft1155.TransactOpts, from, to, id, amount, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.Contract.SetApprovalForAll(&_Nft1155.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.Contract.SetApprovalForAll(&_Nft1155.TransactOpts, operator, approved)
}

// Nft1155ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Nft1155 contract.
type Nft1155ApprovalForAllIterator struct {
	Event *Nft1155ApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Nft1155ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155ApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Nft1155ApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Nft1155ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155ApprovalForAll represents a ApprovalForAll event raised by the Nft1155 contract.
type Nft1155ApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*Nft1155ApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155ApprovalForAllIterator{contract: _Nft1155.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Nft1155ApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155ApprovalForAll)
				if err := _Nft1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) ParseApprovalForAll(log types.Log) (*Nft1155ApprovalForAll, error) {
	event := new(Nft1155ApprovalForAll)
	if err := _Nft1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155TransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the Nft1155 contract.
type Nft1155TransferBatchIterator struct {
	Event *Nft1155TransferBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Nft1155TransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155TransferBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Nft1155TransferBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Nft1155TransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155TransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155TransferBatch represents a TransferBatch event raised by the Nft1155 contract.
type Nft1155TransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*Nft1155TransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155TransferBatchIterator{contract: _Nft1155.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *Nft1155TransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155TransferBatch)
				if err := _Nft1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) ParseTransferBatch(log types.Log) (*Nft1155TransferBatch, error) {
	event := new(Nft1155TransferBatch)
	if err := _Nft1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155TransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the Nft1155 contract.
type Nft1155TransferSingleIterator struct {
	Event *Nft1155TransferSingle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Nft1155TransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155TransferSingle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Nft1155TransferSingle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Nft1155TransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155TransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155TransferSingle represents a TransferSingle event raised by the Nft1155 contract.
type Nft1155TransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*Nft1155TransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155TransferSingleIterator{contract: _Nft1155.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *Nft1155TransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155TransferSingle)
				if err := _Nft1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) ParseTransferSingle(log types.Log) (*Nft1155TransferSingle, error) {
	event := new(Nft1155TransferSingle)
	if err := _Nft1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155URIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the Nft1155 contract.
type Nft1155URIIterator struct {
	Event *Nft1155URI // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Nft1155URIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155URI)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Nft1155URI)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Nft1155URIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155URIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155URI represents a URI event raised by the Nft1155 contract.
type Nft1155URI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*Nft1155URIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155URIIterator{contract: _Nft1155.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) WatchURI(opts *bind.WatchOpts, sink chan<- *Nft1155URI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155URI)
				if err := _Nft1155.contract.UnpackLog(event, "URI", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) ParseURI(log types.Log) (*Nft1155URI, error) {
	event := new(Nft1155URI)
	if err := _Nft1155.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
