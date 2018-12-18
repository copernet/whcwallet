package logic

import (
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"math"
)

func EstimateFee(c *gin.Context) (*view.FeeRate) {
	//Load MemPool feeRate
	feeRate, err := view.GetRPCIns().EstimateFee(2)
	if err == nil {
		return transferRate(feeRate)
	}
	log.WithCtx(c).Errorf("rpc.EstimateFee error:", err.Error())

	return transferRate(0.00001640)
}

func transferRate(feeRate float64) *view.FeeRate {

	fast, _ := decimal.NewFromFloat(feeRate).Mul(decimal.NewFromFloat(1 + config.FeeScale)).Float64()
	slow, _ := decimal.NewFromFloat(feeRate).Mul(decimal.NewFromFloat(1 - config.FeeScale)).Float64()

	return &view.FeeRate{Fast: Round(fast, 8), Normal: Round(feeRate, 8), Slow: Round(slow, 8)}
}

func Round(f float64, n int) float64 {
	res := math.Max(f, 0.00001)
	n10 := math.Pow10(n)
	return math.Trunc((res+0.5/n10)*n10) / n10
}
