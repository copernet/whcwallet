package api

import (
	"strconv"

	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
)

func ListActiveCrowdSales(c *gin.Context) {
	// support property_name and property_id filter
	keyword := c.DefaultQuery("keyword", "")

	total, err := model.GetAllActiveCrowdSaleCount(keyword)
	if err != nil {
		log.WithCtx(c).Errorf("get active crowdsale list total count failed: %v", err)

		c.JSON(200, apiError(ErrGetActiveCrowdSaleCount))
		return
	}

	if total == 0 {
		c.JSON(200, apiSuccess(map[string]interface{}{
			"total": 0,
			"list":  []string{},
		}))
		return
	}

	pageSize, pageNo := paginator(c)
	ret, err := model.GetAllActiveCrowdSale(keyword, pageSize, pageNo)
	if err != nil {
		log.WithCtx(c).Errorf("get active crowd sale list error: %v", err)

		c.JSON(200, apiError(ErrGetActiveCrowdSale))
		return
	}

	crowdSaleList := make([]util.Result, 0, len(ret))
	for _, item := range ret {
		tmp := util.JsonStringToMap(item.PropertyData)
		tmp["created"] = item.BlockTime
		crowdSaleList = append(crowdSaleList, tmp)
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": total,
		"list":  crowdSaleList,
	}))
}

func PurchaseCrowdSaleList(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil || !IsAvailablePropertyID(uint64(pid)) {
		c.JSON(200, apiError(ErrIncorrectPropertyID))
		return
	}

	total, err := model.GetPurchasedCrowdsaleNumber(uint64(pid))
	if err != nil {
		log.WithCtx(c).Errorf("Get purchased crowdsale list total number failed: %v", err)
		c.JSON(200, apiError(ErrGetPurchasedCrowdSaleCount))
		return
	}

	if total == 0 {
		c.JSON(200, apiSuccess(map[string]interface{}{
			"total": 0,
			"list":  []string{},
		}))
		return
	}

	pageSize, pageNo := paginator(c)
	purchaseList, err := model.ListPurchaseCrowdsaleTxes(uint64(pid), pageSize, pageNo)
	if err != nil {
		log.WithCtx(c).Errorf("Get purchased crowdsale list failed: %v", err)
		c.JSON(200, apiError(ErrGetPurchasedCrowdSaleList))
		return
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": total,
		"list":  purchaseList,
	}))
}

func GetPurchasedCrowdSaleTimes(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil || !IsAvailablePropertyID(uint64(pid)) {
		c.JSON(200, apiError(ErrIncorrectPropertyID))
		return
	}

	counts, err := model.GetPurchasedCrowdsaleNumber(uint64(pid))
	if err != nil {
		log.WithCtx(c).Errorf("Get purchased crowdsale count failed: %v", err)
		c.JSON(200, apiError(ErrGetPurchasedCrowdSaleCount))
		return
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": counts,
	}))
}
