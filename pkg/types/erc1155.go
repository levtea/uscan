package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/uchainorg/uscan/pkg/field"
)

type Erc1155Transfer struct {
	TransactionHash common.Hash
	BlockNumber     field.BigInt
	Contract        common.Address
	Method          []byte
	From            common.Address
	To              common.Address
	TokenID         field.BigInt
	Quantity        field.BigInt
	TimeStamp       field.BigInt
}

func (b *Erc1155Transfer) Marshal() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *Erc1155Transfer) Unmarshal(bin []byte) error {
	return rlp.DecodeBytes(bin, &b)
}
