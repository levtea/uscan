package core

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/uchainorg/uscan/pkg/log"
	"github.com/uchainorg/uscan/pkg/storage/fulldb"
	"github.com/uchainorg/uscan/pkg/types"
)

func (n *blockHandle) writeContract(ctx context.Context, data map[common.Address]*types.Contract) (err error) {
	for k, v := range data {
		if err = fulldb.WriteContract(ctx, n.db, k, v); err != nil {
			log.Errorf("write contract(%s): %v ", k, err)
			return err
		}
	}
	return nil
}

func (n *blockHandle) writeProxyContract(ctx context.Context, data map[common.Address]common.Address) (err error) {
	for k, v := range data {
		if err = fulldb.WriteProxyContract(ctx, n.db, k, v); err != nil {
			log.Errorf("write proxy contract(%s): %v ", k, err)
			return err
		}
	}
	return nil
}
