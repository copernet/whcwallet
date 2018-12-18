package api

import (
	"strconv"

	"github.com/bcext/cashutil"
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
)

type ListByOwnerResult struct {
	TxData       util.Result
	PropertyData util.Result
}

func ListProperties(c *gin.Context) {
	// support property_name and property_id filter
	keyword := c.DefaultQuery("keyword", "")

	total, err := model.PropertyListCount(keyword)
	if err != nil {
		log.WithCtx(c).Errorf("get property total number failed: %v", err)
		c.JSON(200, apiError(ErrListPropertiesCount))
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
	list, err := model.ListAllProperties(keyword, pageSize, pageNo)
	if err != nil {
		log.WithCtx(c).Errorf("get properties list failed: %v", err)

		c.JSON(200, apiError(ErrListProperties))
		return
	}

	listProperties := make([]util.Result, 0, len(list))
	for _, item := range list {
		tmp := util.JsonStringToMap(item.PropertyData)
		tmp["created"] = item.BlockTime
		listProperties = append(listProperties, tmp)
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": total,
		"list":  listProperties,
	}))
}

func ListByOwner(c *gin.Context) {
	addresses, err := CheckAddressList(c)
	if err != nil {
		log.WithCtx(c).Errorf("parameter address error: %v", err)
		return
	}

	total, err := model.GetPropertiesNumberByAddresses(addresses)
	if err != nil {
		log.WithCtx(c).Errorf("Get property list issued number by the address failed: %v", err)
		c.JSON(200, apiError(ErrGetPropertyByAddressIssuerCount))
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
	sps, err := model.GetPropertiesByAddresses(addresses, pageSize, pageNo)
	if err != nil {
		log.WithCtx(c).Errorf("get properties from database for address error: %v", err)

		c.JSON(200, apiError(ErrGetPropertyByAddressIssuer))
		return
	}

	result := make([]ListByOwnerResult, len(sps))
	for idx, value := range sps {
		var item ListByOwnerResult

		item.TxData = util.JsonStringToMap(value.TxData)
		item.PropertyData = util.JsonStringToMap(value.PropertyData)

		result[idx] = item
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": total,
		"list":  result,
	}))
}

func GetPropertyByID(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil || !IsAvailablePropertyID(uint64(pid)) {
		c.JSON(200, apiError(ErrIncorrectPropertyID))
		return
	}

	property, err := model.GetPropertyById(uint64(pid))
	if err != nil {
		log.WithCtx(c).Errorf("fetch property data by propertyID failed: %v", err)
		c.JSON(200, apiError(ErrGetPropertyByID))
		return
	}

	c.JSON(200, apiSuccess(util.JsonStringToMap(property)))
}

func IsAvailablePropertyID(id uint64) bool {
	if id < 0 || id == 2 {
		return false
	}

	return true
}

func GetPropertyByName(c *gin.Context) {
	// gin web framework guarantees name is not empty string
	name := c.Param("name")

	if !IsAvailablePropertyName(name) {
		c.JSON(200, apiError(ErrIncorrectPropertyName))
		return
	}

	property, err := model.GetPropertyByName(name)
	if err != nil {
		log.WithCtx(c).Errorf("fetch property data by property_name failed: %v", err)
		c.JSON(200, apiError(ErrGetPropertyByName))
		return
	}

	c.JSON(200, apiSuccess(util.JsonStringArrayToMap(property)))
}

func IsAvailablePropertyName(name string) bool {
	if len(name) == 0 || len(name) > propertyNameLimit {
		return false
	}

	return true
}

// GetProperty query property data via property name of property id.
// will return most ten items for select.
func GetProperty(c *gin.Context) {
	keyword := c.PostForm("keyword")

	if !IsAvailableQueryForProperty(keyword) {
		c.JSON(200, apiError(ErrIncorrectPropertyQueryByKeyword))
		return
	}

	property, err := model.GetPropertyByKeyword(keyword)
	if err != nil {
		log.WithCtx(c).Errorf("fetch property data by keyword(id/name) failed: %v", err)
		c.JSON(200, apiError(ErrGetPropertyByName))
		return
	}

	c.JSON(200, apiSuccess(util.JsonStringArrayToMap(property)))
}

func IsAvailableQueryForProperty(keyword string) bool {
	if len(keyword) > propertyNameLimit {
		return false
	}

	return true
}

func GetPropertyByAddress(c *gin.Context) {
	addr, err := cashutil.DecodeAddress(c.Param("address"), config.GetChainParam())
	if err != nil {
		c.JSON(200, apiError(ErrIncorrectAddress))
		return
	}

	propertyList, err := model.GetPropertyListByAddress(addr.EncodeAddress(true))
	if err != nil {
		log.WithCtx(c).Errorf("Get property list by address failed: %v", err)
		c.JSON(200, apiError(ErrGetPropertyByAddress))
		return
	}

	c.JSON(200, apiSuccess(propertyList))
}

func ListOwners(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("id"))
	total, err := model.ListOwnersCount(pid)
	if err != nil {
		log.WithCtx(c).Errorf("get address balance total number failed: %v", err)
		c.JSON(200, apiError(ErrListPropertiesCount))
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
	list, err := model.ListOwners(pageSize, pageNo, pid)
	if err != nil {
		log.WithCtx(c).Errorf("get properties list failed: %v", err)

		c.JSON(200, apiError(ErrListProperties))
		return
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": total,
		"list":  list,
	}))
}
