package core

import (
	"context"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/uchainorg/uscan/pkg/contract"
	"github.com/uchainorg/uscan/pkg/contract/eip"
	"github.com/uchainorg/uscan/pkg/field"
	"github.com/uchainorg/uscan/pkg/kv"
	"github.com/uchainorg/uscan/pkg/log"
	"github.com/uchainorg/uscan/pkg/storage/forkdb"
	"github.com/uchainorg/uscan/pkg/storage/fulldb"
	"github.com/uchainorg/uscan/pkg/types"
	"github.com/uchainorg/uscan/share"
)

var blockDeleteMap = make(map[*field.BigInt]map[string][][]byte, 0)               // block number => table/key
var blockIndexMap = make(map[*field.BigInt]map[string]*field.BigInt, 0)           // block number => key/index
var blockTotalMap = make(map[*field.BigInt]map[string]*field.BigInt, 0)           // block number => table:key/total
var blockAccountMap = make(map[*field.BigInt]map[string]*field.BigInt, 0)         // block number => account, total tx/total
var blockItxMap = make(map[*field.BigInt]map[string]*field.BigInt, 0)             // block number => account, total itx/total
var blockErc20Map = make(map[*field.BigInt]map[string]*field.BigInt, 0)           // block number => erc20 account, total tx/total
var blockErc721Map = make(map[*field.BigInt]map[string]*field.BigInt, 0)          // block number => erc721 account, total tx/total
var blockErc1155Map = make(map[*field.BigInt]map[string]*field.BigInt, 0)         // block number => erc1155 account, total tx/total
var blockErc20ContractMap = make(map[*field.BigInt]map[string]*field.BigInt, 0)   // block number => erc20 contract account, total tx/total
var blockErc721ContractMap = make(map[*field.BigInt]map[string]*field.BigInt, 0)  // block number => erc721 contract account, total tx/total
var blockErc1155ContractMap = make(map[*field.BigInt]map[string]*field.BigInt, 0) // block number => erc1155 contract account, total tx/total

type blockHandle struct {
	blockData            *types.Block
	transactionData      []*types.Tx
	receiptData          []*types.Rt
	contractOrMemberData map[common.Address]*types.Account
	contractInfoMap      map[common.Address]*types.Contract
	proxyContracts       map[common.Address]common.Address
	internalTxs          map[common.Hash][]*types.InternalTx
	callFrames           map[common.Hash]*types.CallFrame
	contractClient       contract.Contractor
	db                   kv.Database

	newAddrTotal    *field.BigInt
	newErc20Total   *field.BigInt
	newErc721Total  *field.BigInt
	newErc1155Total *field.BigInt
}

func newBlockHandle(
	blockData *types.Block,
	transactionData []*types.Tx,
	receiptData []*types.Rt,
	contractOrMemberData map[common.Address]*types.Account,
	contractInfoMap map[common.Address]*types.Contract,
	proxyContracts map[common.Address]common.Address,
	internalTxs map[common.Hash][]*types.InternalTx,
	callFrames map[common.Hash]*types.CallFrame,
	contractClient contract.Contractor,
	db kv.Database,
) *blockHandle {
	return &blockHandle{
		blockData:            blockData,
		transactionData:      transactionData,
		receiptData:          receiptData,
		contractOrMemberData: contractOrMemberData,
		contractInfoMap:      contractInfoMap,
		proxyContracts:       proxyContracts,
		internalTxs:          internalTxs,
		callFrames:           callFrames,
		contractClient:       contractClient,
		db:                   db,
		newAddrTotal:         field.NewInt(0),
		newErc20Total:        field.NewInt(0),
		newErc721Total:       field.NewInt(0),
		newErc1155Total:      field.NewInt(0),
	}
}

func (n *blockHandle) handleMain(ctx context.Context) (err error) {
	err = fulldb.WriteBlock(ctx, n.db, n.blockData.Number, n.blockData)
	if err != nil {
		log.Errorf("write block : %v, block: %s", err, n.blockData.Number.String())
		return err
	}

	n.newAddrTotal, err = n.checkNewAddr(ctx)
	if err != nil {
		log.Errorf("read acccount to merge: %v", err)
		return err
	}

	//if len(n.contractInfoMap) > 0 {
	//	if err = n.writeContract(ctx, n.contractInfoMap); err != nil {
	//		log.Errorf("write contract: %v", err)
	//		return err
	//	}
	//}
	//if len(n.proxyContracts) > 0 {
	//	if err = n.writeProxyContract(ctx, n.proxyContracts); err != nil {
	//		log.Errorf("write proxy contract: %v", err)
	//		return err
	//	}
	//}

	if len(n.transactionData) > 0 {
		if err = n.writeTxAndRtLog(ctx, n.transactionData, n.receiptData); err != nil {
			log.Errorf("write tx and rt: %v", err)
			return err
		}

		if err = n.writeITx(ctx, n.internalTxs); err != nil {
			log.Errorf("write itxs: %v", err)
			return err
		}

		if err = n.writeTraceTx2(ctx, n.callFrames); err != nil {
			log.Errorf("write callFrames: %v", err)
			return err
		}
	}

	// all account about block write to kv
	if err = n.updateAccounts(ctx); err != nil {
		log.Errorf("write account : %v", err)
		return err
	}

	if err = n.updateHome(ctx); err != nil {
		log.Errorf("write home : %v", err)
		return err
	}

	return nil
}

func (n *blockHandle) handleDeleteFork(ctx context.Context, blockNumber *field.BigInt) (err error) {

	for k, v := range blockDeleteMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				for _, v2 := range v1 {
					_, err = n.db.Get(ctx, v2, &kv.ReadOption{Table: k1})
					if err == nil {
						err = n.db.Del(ctx, v2, &kv.WriteOption{Table: k1})
						if err != nil {
							return err
						}
					}
				}
			}
			delete(blockDeleteMap, k)
		}
	}

	for k, v := range blockIndexMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				i := &field.BigInt{}
				bytesRes, err := n.db.Get(ctx, []byte(k1), &kv.ReadOption{Table: share.ForkIndexTbl})
				if err != nil {
					if errors.Is(err, kv.NotFound) {
						i = field.NewInt(0)
						err = nil
					} else {
						return err
					}
				} else {
					i.SetBytes(bytesRes)
				}
				i.Add(v1)
				err = n.db.Put(ctx, []byte(k1), i.Bytes(), &kv.WriteOption{Table: share.ForkIndexTbl})
				if err != nil {
					return err
				}
			}
			delete(blockIndexMap, k)
		}
	}

	for k, v := range blockTotalMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				i := &field.BigInt{}
				arr := strings.Split(k1, ":")
				tableName := arr[0]
				key := []byte(arr[1])
				bytesRes, err := n.db.Get(ctx, key, &kv.ReadOption{Table: tableName})
				if err != nil {
					return err
				}
				i.SetBytes(bytesRes)
				i.Sub(v1)
				err = n.db.Put(ctx, key, i.Bytes(), &kv.WriteOption{Table: tableName})
				if err != nil {
					return err
				}
			}
			delete(blockTotalMap, k)
		}
	}

	for k, v := range blockAccountMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldAccountTxTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockAccountMap, k)
		}
	}

	for k, v := range blockItxMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldAccountITxTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockItxMap, k)
		}
	}

	for k, v := range blockErc20Map {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldErc20TransferAccountTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockErc20Map, k)
		}
	}

	for k, v := range blockErc721Map {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldErc721TransferAccountTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockErc721Map, k)
		}
	}

	for k, v := range blockErc1155Map {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldErc1155TransferAccountTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockErc1155Map, k)
		}
	}

	for k, v := range blockErc20ContractMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldErc20TransferContractTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockErc20ContractMap, k)
		}
	}

	for k, v := range blockErc721ContractMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldErc721TransferContractTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockErc721ContractMap, k)
		}
	}

	for k, v := range blockErc1155ContractMap {
		if k.Cmp(blockNumber) == 0 {
			for k1, v1 := range v {
				forkOldErc1155TransferContractTotalMap.Add(k1, v1.Bytes())
			}
			delete(blockErc1155ContractMap, k)
		}
	}

	return nil
}

func (n *blockHandle) handleFork(ctx context.Context) (err error) {

	var deleteMap = make(map[string][][]byte, 0)                    // table => key
	var indexMap = make(map[string]*field.BigInt, 0)                // key => index
	var totalMap = make(map[string]*field.BigInt, 0)                // table:key => total
	var accountTotalMap = make(map[string]*field.BigInt, 0)         // account/total tx => total
	var iTxTotalMap = make(map[string]*field.BigInt, 0)             // account/total itx  => total
	var erc20TotalMap = make(map[string]*field.BigInt, 0)           // erc20 transfer tx  => total
	var erc721TotalMap = make(map[string]*field.BigInt, 0)          // erc721 transfer tx  => total
	var erc1155TotalMap = make(map[string]*field.BigInt, 0)         // erc1155 transfer tx  => total
	var erc20ContractTotalMap = make(map[string]*field.BigInt, 0)   // erc20 transfer contract tx  => total
	var erc721ContractTotalMap = make(map[string]*field.BigInt, 0)  // erc20 transfer contract tx  => total
	var erc1155ContractTotalMap = make(map[string]*field.BigInt, 0) // erc20 transfer contract tx  => total

	err = forkdb.WriteBlock(ctx, n.db, n.blockData.Number, n.blockData)
	if err != nil {
		log.Errorf("write fork block : %v, block: %s", err, n.blockData.Number.String())
		return err
	}
	deleteMap[share.ForkBlockTbl] = append(deleteMap[share.ForkBlockTbl], append([]byte("/fork/block/"), n.blockData.Number.Bytes()...))

	n.newAddrTotal, err = n.checkForkNewAddr(ctx)
	if err != nil {
		log.Errorf("read fork account to merge: %v", err)
		return err
	}

	//if len(n.contractInfoMap) > 0 {
	//	if err = n.writeContract(ctx, n.contractInfoMap); err != nil {
	//		log.Errorf("write contract: %v", err)
	//		return err
	//	}
	//}
	//if len(n.proxyContracts) > 0 {
	//	if err = n.writeProxyContract(ctx, n.proxyContracts); err != nil {
	//		log.Errorf("write proxy contract: %v", err)
	//		return err
	//	}
	//}

	if len(n.transactionData) > 0 {
		if err = n.writeForkTxAndRtLog(ctx, n.transactionData, n.receiptData, deleteMap, indexMap, totalMap, accountTotalMap, erc20TotalMap, erc721TotalMap, erc1155TotalMap, erc20ContractTotalMap, erc721ContractTotalMap, erc1155ContractTotalMap); err != nil {
			log.Errorf("write fork tx and rt: %v", err)
			return err
		}

		if err = n.writeForkITx(ctx, n.internalTxs, deleteMap, indexMap, totalMap, iTxTotalMap); err != nil {
			log.Errorf("write fork itxs: %v", err)
			return err
		}

		if err = n.writeForkTraceTx2(ctx, n.callFrames, deleteMap); err != nil {
			log.Errorf("write fork callFrames: %v", err)
			return err
		}
	}

	// all account about block write to kv
	if err = n.updateForkAccounts(ctx); err != nil {
		log.Errorf("write fork account : %v", err)
		return err
	}

	if err = n.updateForkHome(ctx); err != nil {
		log.Errorf("write fork home : %v", err)
		return err
	}

	blockDeleteMap[n.blockData.Number] = deleteMap
	blockIndexMap[n.blockData.Number] = indexMap
	blockTotalMap[n.blockData.Number] = totalMap
	blockAccountMap[n.blockData.Number] = accountTotalMap
	blockItxMap[n.blockData.Number] = iTxTotalMap
	blockErc20Map[n.blockData.Number] = erc20TotalMap
	blockErc721Map[n.blockData.Number] = erc721TotalMap
	blockErc1155Map[n.blockData.Number] = erc1155TotalMap
	blockErc20ContractMap[n.blockData.Number] = erc20ContractTotalMap
	blockErc721ContractMap[n.blockData.Number] = erc721ContractTotalMap
	blockErc1155ContractMap[n.blockData.Number] = erc1155ContractTotalMap

	return nil
}

func (n *blockHandle) handleContractData(ctx context.Context) (err error) {
	if len(n.contractInfoMap) > 0 {
		if err = n.writeContract(ctx, n.contractInfoMap); err != nil {
			log.Errorf("write contract: %v", err)
			return err
		}
	}
	if len(n.proxyContracts) > 0 {
		if err = n.writeProxyContract(ctx, n.proxyContracts); err != nil {
			log.Errorf("write proxy contract: %v", err)
			return err
		}
	}
	return nil
}

func (n *blockHandle) writeTxAndRtLog(ctx context.Context, transactionData []*types.Tx, receiptData []*types.Rt) (err error) {

	for i, v := range transactionData {
		err = fulldb.WriteBlockIndex(ctx, n.db, n.blockData.Number, field.NewInt(int64(i)+1), v.Hash)
		if err != nil {
			log.Errorf("write block index(%d): %v", i, err)
			return err
		}
		if err = n.writeTxAndRt(ctx, v, receiptData[i]); err != nil {
			log.Errorf("writeTxAndRt tx(%s): %v", v.Hash.Hex(), err)
			return err
		}

		var (
			erc20Transfer         *eip.Erc20Transfer
			erc721Transfer        *eip.Ieip721Transfer
			erc1155TransferSignle *eip.Ieip1155TransferSingle
			erc1155TransferBatch  *eip.Ieip1155TransferBatch
		)
		for _, rtLog := range receiptData[i].Logs {
			if len(rtLog.Topics) >= 3 {
				switch rtLog.Topics[0] {
				case contract.TransferEventTopic:
					if len(rtLog.Data) > 0 {
						erc20Transfer, err = n.contractClient.Erc20Transfer(rtLog.Address.Hex(), rtLog.ToEthLog())
						if err == nil {
							if err = n.writeErc20Transfer(ctx, &types.Erc20Transfer{
								TransactionHash: v.Hash,
								BlockNumber:     v.BlockNum,
								Contract:        rtLog.Address,
								Method:          v.Method,
								From:            erc20Transfer.From,
								To:              erc20Transfer.To,
								Amount:          (field.BigInt)(*erc20Transfer.Value),
								TimeStamp:       n.blockData.TimeStamp,
							}); err != nil {
								log.Errorf("write erc20Transfer: %v", err)
								return err
							}
						}
					} else {
						erc721Transfer, err = n.contractClient.Erc721Transfer(rtLog.Address.Hex(), rtLog.ToEthLog())
						if err == nil {
							if err = n.writeErc721Transfer(ctx, &types.Erc721Transfer{
								TransactionHash: v.Hash,
								BlockNumber:     v.BlockNum,
								Contract:        rtLog.Address,
								Method:          v.Method,
								From:            erc721Transfer.From,
								To:              erc721Transfer.To,
								TokenId:         (field.BigInt)(*erc721Transfer.TokenId),
								TimeStamp:       n.blockData.TimeStamp,
							}); err != nil {
								log.Errorf("write erc721Transfer: %v", err)
								return err
							}
						}
					}

				case contract.TransferSingleEventTopic:
					erc1155TransferSignle, err = n.contractClient.Erc1155TransferSingle(rtLog.Address.Hex(), rtLog.ToEthLog())
					if err == nil {
						if err = n.writeErc1155Transfer(ctx, &types.Erc1155Transfer{
							TransactionHash: v.Hash,
							BlockNumber:     v.BlockNum,
							Contract:        rtLog.Address,
							Method:          v.Method,
							From:            erc1155TransferSignle.From,
							To:              erc1155TransferSignle.To,
							TokenID:         (field.BigInt)(*erc1155TransferSignle.Id),
							Quantity:        (field.BigInt)(*erc1155TransferSignle.Value),
							TimeStamp:       n.blockData.TimeStamp,
						}); err != nil {
							log.Errorf("write erc1155Transfer single: %v", err)
							return err
						}
					}
				case contract.TransferBatchEventTopic:
					erc1155TransferBatch, err = n.contractClient.Erc1155TransferBatch(rtLog.Address.Hex(), rtLog.ToEthLog())
					if err == nil {
						for i := range erc1155TransferBatch.Ids {
							if err = n.writeErc1155Transfer(ctx, &types.Erc1155Transfer{
								TransactionHash: v.Hash,
								BlockNumber:     v.BlockNum,
								Contract:        rtLog.Address,
								Method:          v.Method,
								From:            erc1155TransferBatch.From,
								To:              erc1155TransferBatch.To,
								TokenID:         (field.BigInt)(*erc1155TransferBatch.Ids[i]),
								Quantity:        (field.BigInt)(*erc1155TransferBatch.Values[i]),
								TimeStamp:       n.blockData.TimeStamp,
							}); err != nil {
								log.Errorf("write erc1155Transfer single: %v", err)
								return err
							}
						}
					}
				}
			}
		}

		if err = n.updateErc20TrasferTotal(ctx); err != nil {
			log.Errorf("update erc20 transfer total: %v", err)
			return err
		}
		if err = n.updateErc721TrasferTotal(ctx); err != nil {
			log.Errorf("update erc721 transfer total: %v", err)
			return err
		}
		if err = n.updateErc1155TrasferTotal(ctx); err != nil {
			log.Errorf("update erc1155 transfer total: %v", err)
			return err
		}
	}

	return n.writeTxTotal(ctx)
}

func (n *blockHandle) writeTraceTx2(ctx context.Context, callFrames map[common.Hash]*types.CallFrame) (err error) {
	for k, v := range callFrames {
		if err = fulldb.WriteTraceTx2(ctx, n.db, k, &types.TraceTx2{
			Res: v.JsonToString(),
		}); err != nil {
			log.Errorf("write trace tx2: %v", err)
			return err
		}
	}
	return nil
}
