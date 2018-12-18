package api

import (
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
)

type GetNotificationParam struct {
	From      int64    `form:"from" binding:"required,gt=0"`
	To        int64    `form:"to"`
	Addresses []string `form:"address" binding:"required"`
}

func GetNotification(c *gin.Context) {
	var param GetNotificationParam
	err := c.ShouldBind(&param)
	if err != nil {
		c.JSON(200, apiErrorWithMsg(ErrFormItems, err.Error()))
		return
	}

	if len(param.Addresses) > maxRequestAddressList {
		c.JSON(200, apiError(ErrExceedMaxAddressRequestLimit))
		return
	}

	addresses, err := util.ConvToCashAddr(param.Addresses, config.GetChainParam())
	if err != nil {
		c.JSON(200, apiError(ErrIncorrectAddress))
		return
	}

	msgs, err := model.FilterNotification(addresses, param.From, param.To)
	if err != nil {
		c.JSON(200, apiError(ErrGetNotification))
		return
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": len(msgs),
		"list":  msgs,
	}))
}
