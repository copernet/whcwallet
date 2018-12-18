package api

import (
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/api/data"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {

	c.JSON(200, apiSuccess(data.Categories))
}

func GetSubCategories(c *gin.Context) {
	category := c.Query("category")

	subcategory, ok := data.Categories[category]
	if !ok {
		log.WithCtx(c).Error("request not existed category")

		c.JSON(200, apiError(ErrCategoryNotFound))
		return
	}

	c.JSON(200, apiSuccess(subcategory))
}
