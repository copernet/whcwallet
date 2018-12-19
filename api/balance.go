package api

import (
	"github.com/copernet/whccommon/log"
	commodel "github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
)

const (
	defaultUnspentList = 5

	defaultPageSize = 50
)

func GetBalanceForAddress(c *gin.Context) {
	addresses, err := CheckAddressList(c)
	if err != nil {
		return
	}

	ret, err := GetBalanceFromCache(addresses, true)
	if err != nil {
		log.WithCtx(c).Errorf("get balance for address failed: %v", err)
		c.JSON(200, apiError(ErrGetBalanceFromDatabase))
		return
	}

	c.JSON(200, apiSuccess(ret))
}

func GetBalanceFromCache(addresses []string, refresh bool) (map[string][]commodel.BalanceForAddress, error) {
	// Check whether need force to refresh or not.
	// If force to refresh balance, we will get balance via mysql database not redis cache.
	// And store the newest balance information in redis at the same time.
	if refresh {
		bal, err := model.GetBalanceForAddresses(addresses)
		if err != nil {
			return nil, err
		}

		storeBalanceForAddress(bal)

		return bal, nil
	}

	ret, err := model.GetBalanceForAddress(addresses)
	if err != nil {
		return nil, err
	}

	// We should query mysql database if some addresses's balance
	// data not found in redis(because the length of the ret is or
	// equal to addresses's)
	if len(ret) != len(addresses) {
		lackBalanceAddrs := make([]string, 0, len(addresses)-len(ret))
		for _, addr := range addresses {
			if _, ok := ret[addr]; !ok {
				lackBalanceAddrs = append(lackBalanceAddrs, addr)
			}
		}
		remainingRet, err := model.GetBalanceForAddresses(lackBalanceAddrs)
		if err != nil {
			return nil, err
		}

		// If not result returned, we return error instead of return partial
		// balance data.
		if len(remainingRet) == 0 {
			return nil, err
		}

		for addr, bal := range remainingRet {
			ret[addr] = bal
		}

		storeBalanceForAddress(ret)
	}

	// to fill address field
	for addr, item := range ret {
		for idx, _ := range item {
			ret[addr][idx].Address = addr
		}
	}

	return ret, nil
}

func storeBalanceForAddress(bal map[string][]commodel.BalanceForAddress) error {
	return model.StoreBalanceForAddress(bal)
}

type Balance struct {
	Address     string  `json:"address"`
	Confirmed   float64 `json:"confirmed"`
	Unconfirmed float64 `json:"unconfirmed"`
}

func GetBCHBalance(c *gin.Context) {
	legacyAddress, err := CheckAddressListReturnLegacy(c)
	if err != nil {
		return
	}

	ret := make([]Balance, len(legacyAddress))
	for i := 0; i < len(legacyAddress); i++ {
		bal, err := getNode().BlockchainAddressGetBalance(legacyAddress[i])
		if err != nil {
			c.JSON(200, apiError(ErrGetBCHBalance))
			return
		}

		ret[i] = Balance{
			// assume the form post are bech32 encoded addresses
			Address:     util.MustConvSingleCashAddr(legacyAddress[i], config.GetChainParam()),
			Confirmed:   bal.Confirmed.ToBCH(),
			Unconfirmed: bal.Unconfirmed.ToBCH(),
		}
	}

	c.JSON(200, apiSuccess(ret))
}
