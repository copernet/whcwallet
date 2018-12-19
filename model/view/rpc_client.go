package view

import (
	"github.com/copernet/whc.go/rpcclient"
	"github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
)

func GetRPCIns() *rpcclient.Client {
	return model.ConnRpc(config.GetConf().RPC)
}
