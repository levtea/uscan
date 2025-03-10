package core

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/uchainorg/uscan/pkg/field"
	"github.com/uchainorg/uscan/pkg/kv"
	"github.com/uchainorg/uscan/pkg/log"
	"github.com/uchainorg/uscan/pkg/storage/forkdb"
	"github.com/uchainorg/uscan/pkg/types"
	"github.com/uchainorg/uscan/pkg/utils"
	"github.com/uchainorg/uscan/share"
)

var (
	forkAccountItxTotalMap    = utils.NewCache()
	forkOldAccountITxTotalMap = utils.NewCache()
)

func (n *blockHandle) writeForkITx(ctx context.Context, itxmap map[common.Hash][]*types.InternalTx, deleteMap map[string][][]byte, indexMap, totalMap, iTxTotalMap map[string]*field.BigInt) (err error) {
	var itxTotal *field.BigInt

	for k, itxs := range itxmap {
		itxTotal, err = forkdb.ReadITxTotal(ctx, n.db, k)
		if errors.Is(err, kv.NotFound) {
			itxTotal = field.NewInt(0)
		} else {
			log.Errorf("get fork itx total: %v", err)
			return err
		}

		for _, v := range itxs {
			v.TimeStamp = n.blockData.TimeStamp
			if err = forkdb.WriteITx(ctx, n.db, k, itxTotal.Add(field.NewInt(1)), v); err != nil {
				log.Errorf("write fork itx(%s): %v", k.Hex(), err)
				return err
			}
			deleteMap[share.ForkTxTbl] = append(deleteMap[share.ForkTxTbl], append(append([]byte("/fork/iTx/"), k.Bytes()...), append([]byte("/"), itxTotal.Bytes()...)...))
			key2 := append(append([]byte("/fork/iTx/"), k.Bytes()...), []byte("/index")...)
			if indexMap[string(key2)] == nil {
				indexMap[string(key2)] = field.NewInt(0)
			}
			indexMap[string(key2)].Add(field.NewInt(1))

			key := &types.InternalTxKey{
				TransactionHash: v.TransactionHash,
				Index:           *itxTotal,
			}
			if v.From != (common.Address{}) {
				if err = n.writeForkAccountItx(ctx, v.From, key, deleteMap, indexMap, totalMap, iTxTotalMap); err != nil {
					log.Errorf("write fork account(from: %s) Itx: %v", v.From.Hex(), err)
				}
			}

			if v.To != (common.Address{}) {
				if err = n.writeForkAccountItx(ctx, v.To, key, deleteMap, indexMap, totalMap, iTxTotalMap); err != nil {
					log.Errorf("write fork account(to: %s) Itx: %v", v.To.Hex(), err)
				}
			}
		}

		var oldTotal = &field.BigInt{}
		if bytesRes, ok := forkOldAccountITxTotalMap.Get("itxTotal:" + k.String()); ok {
			oldTotal.SetBytes(bytesRes.([]byte))
		} else {
			oldTotal = field.NewInt(0)
		}
		itxTotal.Sub(oldTotal)

		if err = forkdb.WriteItxTotal(ctx, n.db, k, itxTotal); err != nil {
			log.Errorf("write fork itx total: %v", err)
			return err
		}

		key3 := append(append([]byte("/fork/iTx/"), k.Bytes()...), []byte("/total")...)
		if totalMap[share.ForkTxTbl+":"+string(key3)] == nil {
			totalMap[share.ForkTxTbl+":"+string(key3)] = field.NewInt(0)
		}
		totalMap[share.ForkTxTbl+":"+string(key3)].Add(itxTotal)

		itxTotal.Add(oldTotal)
		if iTxTotalMap["itxTotal"+k.String()] == nil {
			iTxTotalMap["itxTotal"+k.String()] = field.NewInt(0)
		}
		iTxTotalMap["itxTotal"+k.String()] = itxTotal

	}
	return nil
}

func (n *blockHandle) writeForkAccountItx(ctx context.Context, addr common.Address, data *types.InternalTxKey, deleteMap map[string][][]byte, indexMap, totalMap, iTxTotalMap map[string]*field.BigInt) (err error) {
	var total = &field.BigInt{}
	if bytesRes, ok := forkAccountItxTotalMap.Get(addr); ok {
		total.SetBytes(bytesRes.([]byte))
	} else {
		total, err = forkdb.ReadAccountITxTotal(ctx, n.db, addr)
		if err != nil {
			if errors.Is(err, kv.NotFound) {
				total = field.NewInt(0)
				err = nil
			} else {
				log.Errorf("get fork account itx total: %v", err)
				return err
			}
		}
	}
	total.Add(field.NewInt(1))
	err = forkdb.WriteAccountITxIndex(ctx, n.db, addr, total, data)
	if err != nil {
		log.Errorf("write fork account itx : %v", err)
		return err
	}
	deleteMap[share.ForkAccountsTbl] = append(deleteMap[share.ForkAccountsTbl], append(append([]byte("/fork/"), addr.Bytes()...), append([]byte("/itx/"), total.Bytes()...)...))
	key := append(append([]byte("/fork/"), addr.Bytes()...), []byte("/itx/index")...)
	if indexMap[string(key)] == nil {
		indexMap[string(key)] = field.NewInt(0)
	}
	indexMap[string(key)].Add(field.NewInt(1))

	var oldTotal = &field.BigInt{}
	if bytesRes, ok := forkOldAccountITxTotalMap.Get(addr.String()); ok {
		oldTotal.SetBytes(bytesRes.([]byte))
	} else {
		oldTotal = field.NewInt(0)
	}
	total.Sub(oldTotal)

	err = forkdb.WriteAccountITxTotal(ctx, n.db, addr, total)

	key2 := append(append([]byte("/fork/"), addr.Bytes()...), []byte("/itx/total")...)
	if totalMap[share.ForkAccountsTbl+":"+string(key2)] == nil {
		totalMap[share.ForkAccountsTbl+":"+string(key2)] = field.NewInt(0)
	}
	totalMap[share.ForkAccountsTbl+":"+string(key2)].Add(field.NewInt(1))

	total.Add(oldTotal)
	if iTxTotalMap[addr.String()] == nil {
		iTxTotalMap[addr.String()] = field.NewInt(0)
	}
	iTxTotalMap[addr.String()] = total

	if err == nil {
		forkAccountItxTotalMap.Add(addr, total.Bytes())
	}

	return
}
