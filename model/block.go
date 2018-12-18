package model

import (
	common "github.com/copernet/whccommon/model"
	"github.com/jinzhu/gorm"
)

func GetLastBlock() *common.Block {
	var block = common.Block{}
	err := db.Last(&block).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil
	}

	return &block
}
