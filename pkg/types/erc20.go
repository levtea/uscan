package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/uchainorg/uscan/pkg/field"
)

type Erc20Transfer struct {
	TransactionHash common.Hash

	BlockNumber field.BigInt
	Contract    common.Address
	Method      []byte
	From        common.Address
	To          common.Address
	Amount      field.BigInt
	TimeStamp   field.BigInt
}

func (b *Erc20Transfer) Marshal() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *Erc20Transfer) Unmarshal(bin []byte) error {
	return rlp.DecodeBytes(bin, &b)
}
