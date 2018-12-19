package api

import (
	"github.com/copernet/whcwallet/config"
	"github.com/gin-gonic/gin"
)

func GetEnv(c *gin.Context) {
	if config.GetConf().TestNet {
		c.JSON(200, apiSuccess(map[string]bool{"testnet": true}))
		return
	}

	c.JSON(200, apiSuccess(map[string]bool{"testnet": false}))
	return
}
