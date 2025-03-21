/*
Copyright © 2022 uscan team

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package pkg

import (
	"context"
	"os"

	"github.com/uchainorg/uscan/pkg/service"
	"github.com/uchainorg/uscan/pkg/storage"

	"github.com/spf13/viper"
	"github.com/uchainorg/uscan/pkg/contract"
	"github.com/uchainorg/uscan/pkg/core"
	"github.com/uchainorg/uscan/pkg/rpcclient"
	"github.com/uchainorg/uscan/share"

	"github.com/spf13/cobra"
	"github.com/sunvim/utils/grace"
	"github.com/uchainorg/uscan/pkg/apis"
)

func MainRun(cmd *cobra.Command, args []string) {
	os.RemoveAll(viper.GetString(share.MdbxPath) + "/fork")
	storage := storage.NewStorage(viper.GetString(share.MdbxPath))
	rpcMgr := rpcclient.NewRpcClient(viper.GetStringSlice(share.RpcUrls))

	sync := core.NewSync(rpcMgr, contract.NewClient(rpcMgr), viper.GetInt64(share.ForkBlockNum), storage.FullDB, storage.ForkDB, viper.GetUint64(share.WorkChan))
	go sync.Execute(context.Background())

	service.NewStore(storage)
	service.StartHandleContractVerity()
	apis.GetChainID(rpcMgr.ChainID(context.Background()))
	_, svc := grace.New(context.Background())
	svc.RegisterService("web service", apis.Apis)
	svc.Wait()
}
