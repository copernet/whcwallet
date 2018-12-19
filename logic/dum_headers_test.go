package logic

import (
	"fmt"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/whccommon/log"
	"github.com/qshuai/tcolor"
	"os"
	"testing"
)

func TestDumpHeaders(t *testing.T) {
	DumpHeaders(log.NewContext(), nil)
}

func TestLoadHeaderFile(t *testing.T) {

	path := ""
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		log.WithCtx(nil).Errorf("OpenFile.error:%s", err.Error())
		return
	}
	fi, err := file.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		log.WithCtx(nil).Errorf("OpenFile.error:%s", err.Error())
		return
	}

	size := fi.Size() / 80
	for i := size-1; i >=0; i-- {
		_, err := file.Seek(i*80, 0)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "file seek failed"))
			return
		}

		var blockheader wire.BlockHeader
		err = blockheader.Deserialize(file)
		if err != nil {
			log.WithCtx(nil).Errorf("Deserialize.error:%s", err.Error())
			return
		}

		fmt.Println(blockheader.BlockHash())
	}

}
