package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/uchainorg/uscan/pkg/field"
)

type Log struct {
	Address  common.Address `json:"address"`
	Topics   []common.Hash  `json:"topics"`
	Data     hexutil.Bytes  `json:"data"`
	LogIndex field.BigInt   `json:"logIndex"`
}

func (l *Log) ToEthLog() ethTypes.Log {
	return ethTypes.Log{
		Address: l.Address,
		Topics:  l.Topics,
		Data:    l.Data,
	}
}

type Rt struct {
	TxHash            common.Hash     `json:"transactionHash"  rlp:"-"`
	Type              field.BigInt    `json:"type,omitempty"`
	PostState         hexutil.Bytes   `json:"root"`
	Status            field.BigInt    `json:"status"`
	CumulativeGasUsed field.BigInt    `json:"cumulativeGasUsed"`
	Bloom             Bloom           `json:"logsBloom"`
	Logs              []*Log          `json:"logs"`
	ContractAddress   *common.Address `json:"contractAddress"`
	GasUsed           field.BigInt    `json:"gasUsed"`
	EffectiveGasPrice field.BigInt    `json:"effectiveGasPrice"`
	ExistInternalTx   bool
	ReturnErr         string
}

func (b *Rt) Marshal() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *Rt) Unmarshal(bin []byte) error {
	return rlp.DecodeBytes(bin, &b)
}
