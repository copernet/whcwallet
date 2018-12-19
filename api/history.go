package api

import (
	"net/http"
	"strconv"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/copernet/whccommon/log"
	common "github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/model/view"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

var historyCountCache util.CacheMap

const (
	bitcoinCashDetail = "bch"
)

func init() {
	historyCountCache.New()
}

func GetHistory(c *gin.Context) {
	result := make(map[string]interface{})

	s_id := c.Param("id")
	ppId, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		c.JSON(200, apiError(ErrConvertInt))
		return
	}

	address := c.PostForm("address")
	if address != "" {
		addr, err := cashutil.DecodeAddress(address, &chaincfg.TestNet3Params)
		if err != nil {
			c.JSON(200, apiError(ErrIncorrectAddress))
			return
		}

		address = addr.EncodeAddress(true)
	}

	sStart := c.PostForm("start")
	start, err := strconv.ParseInt(sStart, 10, 64)
	if err != nil {
		c.JSON(200, apiError(ErrConvertInt))
		return
	}

	sCount := c.PostForm("count")
	count, err := strconv.ParseInt(sCount, 10, 64)
	if err != nil {
		c.JSON(200, apiError(ErrConvertInt))
		return
	}

	total, err := model.QueryTotal(ppId, address)
	if err != nil {
		c.JSON(200, apiError(ErrQueryTotal))
		return
	}

	result["total"] = total

	transactions, err := model.QueryTransactions(ppId, address, start, count)
	if err != nil {
		log.WithCtx(c).Errorf("query transaction error from database: %v", err)

		c.JSON(200, apiError(ErrQueryTransactions))
		return
	}

	list := util.JsonStringArrayToMap(transactions)
	for idx, item := range list {
		v, ok := item["amount"]
		if ok {
			if str, ok := v.(string); ok {
				amount, err := decimal.NewFromString(str)
				if err != nil {
					log.WithCtx(c).Errorf("Fatal error! Server database stored amount string malformed: %s", str)

					c.JSON(200, apiError(ErrAmountString))
					return
				}

				list[idx]["amount"] = amount.String()
			}
		}
	}

	result["transactions"] = list

	c.JSON(http.StatusOK, apiSuccess(result))
}

func GetHistoryDetail(c *gin.Context) {
	tx_hash := c.Query("tx_hash")
	if len(tx_hash) != 64 {
		c.JSON(200, apiError(ErrHash256Format))
		return
	}

	isBitcoinCash := c.Query("type")
	if isBitcoinCash == bitcoinCashDetail {
		txDetail, err := getNode().BlockchainTransactionGet(tx_hash, true)
		if err != nil {
			log.WithCtx(c).Info("get transaction via electrum server failed")
			c.JSON(200, apiError(ErrTransactionGet))
			return
		}

		c.JSON(200, apiSuccess(txDetail))
		return
	}

	detail, err := model.GetHistoryDetail(tx_hash)
	if err != nil {
		log.WithCtx(c).Errorf("get history detail from database error: %v", err)

		c.JSON(200, apiError(ErrGetHistoryDetail))
		return
	}

	c.JSON(200, apiSuccess(util.JsonStringToMap(detail)))
}

func GetHistoryDetailPending(c *gin.Context) {
	tx_hash := c.Query("tx_hash")
	if len(tx_hash) != 64 {
		c.JSON(200, apiError(ErrHash256Format))
		return
	}

	cient := view.GetRPCIns()
	ret, err := cient.WhcGetTransaction(tx_hash)
	if err != nil {
		log.WithCtx(c).Errorf("get pending wormhole transaction failed: %v", err)
		c.JSON(200, apiError(ErrWhcGetTransaction))
		return
	}

	c.JSON(200, apiSuccess(ret))
}

type HistoryListResult struct {
	TxHash                      string `json:"tx_hash"`
	TxType                      int32  `json:"tx_type"`
	TxState                     string `json:"tx_state"`
	Address                     string `json:"address"`
	AddressRole                 string `json:"address_role"`
	BalanceAvailableCreditDebit string `json:"balance_available_credit_debit"`
	PropertyID                  uint64 `json:"property_id"`
	PropertyName                string `json:"property_name"`
	Created                     int64  `json:"created"`
}

func GetHistoryList(c *gin.Context) {
	addresses, err := CheckAddressList(c)
	if err != nil {
		log.WithCtx(c).Errorf("parameter address error: %v", err)
		return
	}

	pageSize, pageNo := paginator(c)
	var pid int
	propertyID, ok := c.GetPostForm("property_id")
	if ok {
		pid, err = strconv.Atoi(propertyID)
		if err != nil {
			c.JSON(200, apiError(ErrIncorrectPropertyID))
			return
		}

		if pid == 0 {
			list, total, err := getBCHHistory(addresses[0], pageNo, pageSize)
			if err != nil {
				c.JSON(200, apiError(ErrGetHistoryList))
				return
			}

			c.JSON(200, apiSuccess(map[string]interface{}{
				"total": total,
				"list":  formatHistoryList(list),
			}))
			return
		}
	}

	total, err := model.GetHistoryListCount(addresses, pid)
	if err != nil {
		log.WithCtx(c).Errorf("get history list total count failed: %v", err)
		c.JSON(200, apiError(ErrGetHistoryListCount))
		return
	}

	if total == 0 {
		c.JSON(200, apiSuccess(map[string]interface{}{
			"total": 0,
			"list":  []string{},
		}))
		return
	}

	list, err := model.GetHistoryList(addresses, pid, pageSize, pageNo)
	if err != nil {
		log.WithCtx(c).Warnf("get transaction history list error: %v", err)

		c.JSON(200, apiError(ErrGetHistoryList))
		return
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"total": total,
		"list":  formatHistoryList(list),
	}))
}

func formatHistoryList(list []model.HistoryList) []HistoryListResult {
	var resultSet []HistoryListResult
	for _, v := range list {
		var temp HistoryListResult
		temp.TxState = v.TxState
		temp.TxHash = v.TxHash
		temp.TxType = v.TxType
		temp.Address = v.Address
		temp.AddressRole = v.AddressRole
		temp.BalanceAvailableCreditDebit = v.BalanceAvailableCreditDebit.String()
		temp.PropertyName = v.PropertyName
		temp.PropertyID = v.PropertyID
		// reset BlockTime to zero for web display `unconfirmed`
		if v.TxState == string(common.Pending) {
			temp.Created = 0
		} else {
			temp.Created = v.BlockTime
		}
		resultSet = append(resultSet, temp)
	}

	return resultSet
}
