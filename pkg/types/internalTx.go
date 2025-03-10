package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/uchainorg/uscan/pkg/field"
)

type InternalTx struct {
	TransactionHash common.Hash `rlp:"-"`
	BlockNumber     field.BigInt
	Status          bool
	CallType        string
	Depth           string
	From            common.Address
	To              common.Address
	Amount          field.BigInt
	GasLimit        field.BigInt
	TimeStamp       field.BigInt
}

func (b *InternalTx) Marshal() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *InternalTx) Unmarshal(bin []byte) error {
	return rlp.DecodeBytes(bin, &b)
}

type InternalTxKey struct {
	TransactionHash common.Hash
	Index           field.BigInt
}

func (b *InternalTxKey) Marshal() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *InternalTxKey) Unmarshal(bin []byte) error {
	return rlp.DecodeBytes(bin, &b)
}
