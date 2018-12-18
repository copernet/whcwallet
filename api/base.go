package api

import (
	"errors"
	"strconv"

	"github.com/bcext/cashutil"
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
)

const (
	propertyNameLimit = 520

	// pagination feature
	defaultPageSizeNumber = 50
	defaultPageNo         = 1

	// user request limit
	maxRequestAddressList = 50

	// user request address field name
	addressListFieldName = "address"
)

var (
	supportPrecision map[string]struct{}
)

type Response struct {
	Code    errCode     `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func IsSupportPrecision(precision string) bool {
	if _, ok := supportPrecision[precision]; !ok {
		return false
	}

	return true
}

func InTxTypeArray(target int32, array []int32) bool {
	for _, item := range array {
		if item == target {
			return true
		}
	}

	return false
}

// CheckAddressList check address number limit and encoded,
// and only apply for post request.
func CheckAddressList(c *gin.Context) ([]string, error) {
	addresses := c.PostFormArray(addressListFieldName)

	if len(addresses) == 0 {
		c.JSON(200, apiError(ErrEmptyAddressList))
		return nil, errors.New("empty address list")
	}

	if len(addresses) > maxRequestAddressList {
		c.JSON(200, apiError(ErrExceedMaxAddressRequestLimit))
		return nil, errors.New("so many addresses requested")
	}

	// decode all bitcoin cash addresses, return immediately if encounter any error
	addresses, err := util.ConvToCashAddr(addresses, config.GetChainParam())
	if err != nil {
		c.JSON(200, apiError(ErrIncorrectAddress))
		return nil, err
	}

	return addresses, nil
}

func CheckAddressListReturnLegacy(c *gin.Context) ([]string, error) {
	addresses := c.PostFormArray(addressListFieldName)

	if len(addresses) > maxRequestAddressList {
		c.JSON(200, apiError(ErrExceedMaxAddressRequestLimit))
		return nil, errors.New("so many addresses requested")
	}

	param := config.GetChainParam()
	for idx, addr := range addresses {
		// to guarantee all addresses is valid
		cashAddr, err := cashutil.DecodeAddress(addr, param)
		if err != nil {
			c.JSON(200, apiError(ErrIncorrectAddress))
			return nil, err
		}

		addresses[idx] = cashAddr.EncodeAddress(false)
	}

	return addresses, nil
}

func paginator(c *gin.Context) (int, int) {
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(defaultPageSizeNumber)))
	pageNo, err := strconv.Atoi(c.DefaultQuery("pageNo", strconv.Itoa(defaultPageNo)))
	if err != nil {
		log.WithCtx(c).Errorf("query property history list parameter error, pagesize: %s pageno: %s",
			c.Query("pageSize"), c.Query("pageNo"))
	}

	if pageSize <= 0 {
		pageSize = defaultPageSizeNumber
	}

	if pageNo <= 0 {
		pageNo = defaultPageNo
	}

	return pageSize, pageNo
}

func init() {
	supportPrecision = make(map[string]struct{})
	for i := 0; i <= 8; i++ {
		supportPrecision[strconv.Itoa(i)] = struct{}{}
	}

	// initial redis instance
	// get configuration for this project
	err := model.ConnRedis()
	if err != nil {
		panic(err)
	}
}
