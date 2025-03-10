package core

import (
	"context"
	"errors"
	"time"

	"github.com/uchainorg/uscan/pkg/field"
	"github.com/uchainorg/uscan/pkg/kv"
	"github.com/uchainorg/uscan/pkg/log"
	"github.com/uchainorg/uscan/pkg/storage/forkdb"
	"github.com/uchainorg/uscan/pkg/types"
)

var forkHomeCache *types.Home

func (n *blockHandle) updateForkHome(ctx context.Context) (err error) {
	var home *types.Home
	if forkHomeCache == nil {
		home, err = forkdb.ReadHome(ctx, n.db)
		if err != nil {
			if errors.Is(err, kv.NotFound) {
				log.Infof("read fork home not found")
				home = &types.Home{
					TxTotal:      *field.NewInt(0),
					AddressTotal: *field.NewInt(0),
					Blocks:       make([]*types.BkSim, 0, 1),
					Txs:          make([]*types.TxSim, 0, 10),
					DateTxs:      make(map[string]*field.BigInt),
				}
				err = nil
				forkHomeCache = home
			} else {
				return err
			}
		} else {
			forkHomeCache = home
		}
	} else {
		home = forkHomeCache
	}

	home.BlockNumber.SetBytes(n.blockData.Number.Bytes())

	home.TxTotal.Add(field.NewInt(int64(len(n.transactionData))))
	home.AddressTotal.Add(n.newAddrTotal)
	home.Erc20Total.Add(n.newErc20Total)
	home.Erc721Total.Add(n.newErc721Total)
	home.Erc721Total.Add(n.newErc1155Total)
	home.Blocks = append(home.Blocks, &types.BkSim{
		Number:            *n.blockData.Number,
		Timestamp:         n.blockData.TimeStamp,
		Miner:             n.blockData.Coinbase,
		GasUsed:           n.blockData.GasUsed,
		TransactionsTotal: *field.NewInt(int64(len(n.transactionData))),
	})
	for _, v := range n.transactionData {
		home.Txs = append(home.Txs, &types.TxSim{
			Hash:      v.Hash,
			From:      v.From,
			To:        *v.To,
			GasPrice:  v.GasPrice,
			Gas:       v.Gas,
			Timestamp: v.TimeStamp,
		})
	}

	timeLayout := "20060102"
	timeStr := time.Unix(int64(n.blockData.TimeStamp.ToUint64()), 0).UTC().Format(timeLayout)
	if _, ok := home.DateTxs[timeStr]; ok {
		home.DateTxs[timeStr].Add(field.NewInt(int64(len(n.transactionData))))
	} else {
		home.DateTxs[timeStr] = field.NewInt(int64(len(n.transactionData)))
	}

	if len(home.Blocks) > 10 {
		home.Blocks = home.Blocks[(len(home.Blocks) - 10):]
	}

	if len(home.Txs) > 10 {
		home.Txs = home.Txs[(len(home.Txs) - 10):]
	}

	delete(home.DateTxs, time.Unix(int64(n.blockData.TimeStamp.ToUint64()-(3600*24*14)), 0).UTC().Format(timeLayout))

	if err = forkdb.WriteSyncingBlock(ctx, n.db, n.blockData.Number); err != nil {
		log.Errorf("write fork syncing block: %v", err)
		return err
	}

	return forkdb.WriteHome(ctx, n.db, home)
}
