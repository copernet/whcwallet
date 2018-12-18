package model

import (
	"github.com/copernet/whccommon/model"
)

func GetNewestVersionFroDevice(device string) (*model.Version, error) {
	var version model.Version
	err := db.Table("versions").Where("device = ?", device).First(&version).Error
	return &version, err
}
