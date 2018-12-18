package api

import (
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func UpdateSoftware(c *gin.Context) {
	device := c.Param("device")

	if device == "" {
		c.JSON(200, apiError(ErrEmptyDeviceName))
		return
	}

	v, err := model.GetNewestVersionFroDevice(device)
	if gorm.IsRecordNotFoundError(err) {
		c.JSON(200, &Response{
			Code:    0,
			Message: "It is already the latest version",
		})
		return
	}

	if err != nil {
		c.JSON(200, apiError(ErrGetNewestVersionFailed))
		return
	}

	ver := view.Version{Version: v.Version, VersionCode: v.VersionCode, Download: v.Download, Description: v.Description}
	c.JSON(200, apiSuccess(ver))
}
