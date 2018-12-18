package logic

import (
	"context"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/go-electrum/electrum"
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/model/view"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultFile = "headers.dat"
)

func DumpHeaders(ctx context.Context, h *electrum.BlockchainHeader) {
	client := view.GetRPCIns()

	start := int64(0)
	end := model.GetLastBlock().BlockHeight

	file, err := os.Open(getBasePath())
	if !os.IsNotExist(err) {
		//For block init
		header, err := loadCurrHeader(file)
		if err == nil {
			blockHash := header.BlockHash()
			block, err := client.GetBlockVerbose(&blockHash)
			if err != nil {
				// Could not obtain stat, handle error
				log.WithCtx(ctx).Errorf("GetBlockVerbose.error:%s", err.Error())
				return
			}

			start = block.Height + 1
			if h != nil && int64(h.BlockHeight) > block.Height {
				start = int64(h.BlockHeight)
			}
		}
	}

	for i := start; i <= end; i++ {
		loadBlockHeader(i, ctx)
		log.WithCtx(ctx).Infof("header for height:%d ok", i)
	}
}
func loadCurrHeader(file *os.File) (*wire.BlockHeader, error) {

	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		return nil, err
	}

	_, err = file.Seek(fi.Size()-80, 0)
	if err != nil {
		return nil, err
	}

	var header wire.BlockHeader
	err = header.Deserialize(file)
	if err != nil {
		return nil, err
	}

	return &header, nil
}
func getBasePath() string {
	path, _ := filepath.Abs("./")
	lastIndex := strings.Index(path, config.ProjectLastDir) + len(config.ProjectLastDir)
	return path[:lastIndex] + "/static/" + DefaultFile

}
func loadBlockHeader(height int64, ctx context.Context) {
	client := view.GetRPCIns()
	blockHash, err := client.GetBlockHash(height)
	if err != nil {
		log.WithCtx(ctx).Errorf("GetBlockHash.error:%s", err.Error())
		return
	}

	header, err := client.GetBlockHeader(blockHash)
	if err != nil {
		log.WithCtx(ctx).Errorf("GetBlockHeader.error:%s", err.Error())
		return
	}

	file, err := os.OpenFile(getBasePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644))
	if os.IsNotExist(err) {
		log.WithCtx(ctx).Errorf("OpenFile.error:%s", err.Error())
	}

	err = header.Serialize(file)
	if err != nil {
		log.WithCtx(ctx).Errorf("Serialize.error:%s", err.Error())
		return
	}
}
