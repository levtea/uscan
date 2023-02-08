package fulldb

import (
	"context"

	"github.com/uchainorg/uscan/pkg/field"
	"github.com/uchainorg/uscan/pkg/kv"
	"github.com/uchainorg/uscan/pkg/types"
	"github.com/uchainorg/uscan/share"
)

var (
	homeKey    = []byte("/home")
	syncingKey = []byte("/syncing")
)

/*
table: home

/home => home
/syncing => block number
*/

func ReadHome(ctx context.Context, db kv.Reader) (home *types.Home, err error) {
	var bytesRes []byte
	bytesRes, err = db.Get(ctx, homeKey, &kv.ReadOption{Table: share.HomeTbl})
	if err != nil {
		return
	}
	home = &types.Home{}
	err = home.Unmarshal(bytesRes)
	return
}

func WriteHome(ctx context.Context, db kv.Writer, home *types.Home) (err error) {
	var bytesRes []byte
	bytesRes, err = home.Marshal()
	if err != nil {
		return
	}
	return db.Put(ctx, homeKey, bytesRes, &kv.WriteOption{Table: share.HomeTbl})
}

func ReadSyncingBlock(ctx context.Context, db kv.Reader) (bk *field.BigInt, err error) {
	var bytesRes []byte
	bytesRes, err = db.Get(ctx, syncingKey, &kv.ReadOption{Table: share.HomeTbl})

	if err != nil {
		return
	}
	bk = &field.BigInt{}
	bk.SetBytes(bytesRes)
	return
}

func WriteSyncingBlock(ctx context.Context, db kv.Writer, bk *field.BigInt) (err error) {
	return db.Put(ctx, syncingKey, bk.Bytes(), &kv.WriteOption{Table: share.HomeTbl})
}
