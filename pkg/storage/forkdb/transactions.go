package forkdb

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/uchainorg/uscan/pkg/field"
	"github.com/uchainorg/uscan/pkg/kv"
	"github.com/uchainorg/uscan/pkg/types"
	"github.com/uchainorg/uscan/share"
)

var (
	txKey      = []byte("/fork/tx/")
	rtKey      = []byte("/fork/rt/")
	txTotalKey = []byte("/fork/all/tx/total")
	txIndexKey = []byte("/fork/all/tx/")
)

/*
table : transactions

/fork/tx/<txhash> => tx info
/fork/rt/<txhash> => rt info
/fork/all/tx/total => total
/fork/all/tx/<index> => <txhash>
*/

func WriteTx(ctx context.Context, db kv.Writer, hash common.Hash, data *types.Tx) (err error) {
	var (
		key      = append(txKey, hash.Bytes()...)
		bytesRes []byte
	)
	bytesRes, err = data.Marshal()
	if err != nil {
		return err
	}
	return db.Put(ctx, key, bytesRes, &kv.WriteOption{Table: share.ForkTxTbl})
}

func ReadTx(ctx context.Context, db kv.Reader, hash common.Hash) (data *types.Tx, err error) {
	var (
		key      = append(txKey, hash.Bytes()...)
		bytesRes []byte
	)

	bytesRes, err = db.Get(ctx, key, &kv.ReadOption{Table: share.ForkTxTbl})
	if err != nil {
		return
	}
	data = &types.Tx{}
	err = data.Unmarshal(bytesRes)
	if err == nil {
		data.Hash = hash
	}
	return
}

func DeleteTx(ctx context.Context, db kv.Writer, hash common.Hash) (err error) {
	var key = append(txKey, hash.Bytes()...)
	return db.Del(ctx, key, &kv.WriteOption{Table: share.ForkTxTbl})
}

func WriteTxIndex(ctx context.Context, db kv.Writer, index *field.BigInt, hash common.Hash) error {
	return db.Put(ctx, append(txIndexKey, index.Bytes()...), hash.Bytes(), &kv.WriteOption{Table: share.ForkTxTbl})
}

func ReadTxByIndex(ctx context.Context, db kv.Reader, index *field.BigInt) (data *types.Tx, err error) {
	var hashByte []byte
	hashByte, err = db.Get(ctx, append(txIndexKey, index.Bytes()...), &kv.ReadOption{Table: share.ForkTxTbl})
	if err != nil {
		return
	}
	hash := common.BytesToHash(hashByte)
	return ReadTx(ctx, db, hash)
}

func WriteTxTotal(ctx context.Context, db kv.Writer, total *field.BigInt) error {
	return db.Put(ctx, txTotalKey, total.Bytes(), &kv.WriteOption{Table: share.ForkTxTbl})
}

func ReadTxTotal(ctx context.Context, db kv.Reader) (total *field.BigInt, err error) {
	var bytesRes []byte
	bytesRes, err = db.Get(ctx, txTotalKey, &kv.ReadOption{Table: share.ForkTxTbl})
	if err != nil {
		return
	}
	total = &field.BigInt{}
	total.SetBytes(bytesRes)
	return
}

func WriteRt(ctx context.Context, db kv.Writer, hash common.Hash, data *types.Rt) (err error) {
	var (
		key      = append(rtKey, hash.Bytes()...)
		bytesRes []byte
	)
	bytesRes, err = data.Marshal()
	if err != nil {
		return
	}
	return db.Put(ctx, key, bytesRes, &kv.WriteOption{Table: share.ForkTxTbl})
}

func ReadRt(ctx context.Context, db kv.Reader, hash common.Hash) (data *types.Rt, err error) {
	var (
		key      = append(rtKey, hash.Bytes()...)
		bytesRes []byte
	)
	bytesRes, err = db.Get(ctx, key, &kv.ReadOption{Table: share.ForkTxTbl})
	if err != nil {
		return
	}
	data = &types.Rt{}
	err = data.Unmarshal(bytesRes)
	if err == nil {
		data.TxHash = hash
	}
	return
}

func DeleteRt(ctx context.Context, db kv.Writer, hash common.Hash) (err error) {
	var key = append(rtKey, hash.Bytes()...)
	return db.Del(ctx, key, &kv.WriteOption{Table: share.ForkTxTbl})
}
