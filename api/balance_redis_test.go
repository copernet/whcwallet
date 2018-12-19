package api

import (
	"testing"

	"github.com/copernet/whccommon/model"
	"github.com/shopspring/decimal"
)

func TestStoreBalance(t *testing.T) {
	amount1 := decimal.New(232, 0)
	amount2 := decimal.New(232, 0)
	bal := []model.BalanceForAddress{
		{
			Address:          "bchtest:qpncl685d5eys3tj4vykrhw4507fcuas9scm4qd70e",
			PropertyID:       1,
			BalanceAvailable: &amount1,
			//PendingPos:       9900,
			//PendingNeg:       -68,
		},
		{
			Address:          "bchtest:qpncl685d5eys3tj4vykrhw450734uas9scm4qd7od",
			PropertyID:       1,
			BalanceAvailable: &amount2,
			//PendingPos:       9900,
			//PendingNeg:       -68,
		},
	}

	err := storeBalanceForAddress(map[string][]model.BalanceForAddress{
		bal[0].Address: {bal[0]},
		bal[1].Address: {bal[1]},
	})

	if err != nil {
		t.Error(err)
	}
}
