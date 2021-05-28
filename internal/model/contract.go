package model

type ContractSchemaType string

const (
	ContractSchemaTypeERC1155 ContractSchemaType = "ERC1155"
	ContractSchemaTypeERC20 ContractSchemaType = "ERC20"
)

func (t ContractSchemaType) String() string {
	return string(t)
}
