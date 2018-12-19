package routers

import (
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/api"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(ginLogger())
	r.Use(log.LogContext([]string{"/static","/headers/file"}))

	// acquire the current environment (mainnet / testnet)
	r.GET("/env", api.GetEnv)

	// define account apis
	account := r.Group("/user")
	{
		account.GET("/wallet/challenge", api.Challenge)
		account.POST("/wallet/create", api.Create)
		account.POST("/wallet/update", api.Update)
		account.POST("/wallet/login", api.Login)
		account.POST("/wallet/verify", api.Verify)
		account.POST("/wallet/newMfa", api.NewMfa)
	}

	tx := r.Group("/getunsigned")
	tx.Use(IsSupportTxType())
	tx.Use(FixedFormParam())
	{
		tx.POST("/:txtype", api.GetUnsignedTx)
	}

	transaction := r.Group("/transaction")
	{
		transaction.POST("/push", api.PushTx)
		transaction.POST("/fee", api.FeeRate)
	}

	history := r.Group("/history")
	{
		history.POST("/list", api.GetHistoryList)
		history.GET("/detail", api.GetHistoryDetail)
		history.GET("/detail/pending", api.GetHistoryDetailPending)
		history.POST("/id/:id", api.GetHistory)
	}

	property := r.Group("/property")
	{
		property.POST("/listbyowner", api.ListByOwner)
		property.POST("/listowners/:id", api.ListOwners)
		property.GET("/list", api.ListProperties)
		property.GET("/id/:id", api.GetPropertyByID)
		property.GET("/name/:name", api.GetPropertyByName)
		// support fuzzy search via property name or property id
		property.POST("/query", api.GetProperty)
		property.GET("/address/:address", api.GetPropertyByAddress)
	}

	category := r.Group("/category")
	{
		category.GET("/", api.GetCategories)
		category.GET("/subcategories", api.GetSubCategories)
	}

	class := r.Group("/classification")
	{
		class.GET("/", api.GetCategories)
		class.GET("/subcategories", api.GetSubCategories)
	}

	crowdSale := r.Group("/crowdsale")
	{
		// /crowdsale/list/active?pagesize=10&pageno=1
		crowdSale.GET("/list/active", api.ListActiveCrowdSales)
		crowdSale.GET("/purchase/record/id/:id", api.PurchaseCrowdSaleList)
		crowdSale.GET("/purchase/times/id/:id", api.GetPurchasedCrowdSaleTimes)
	}

	balance := r.Group("/balance")
	{
		balance.POST("/addresses", api.GetBalanceForAddress)
		balance.POST("/bch/addresses", api.GetBCHBalance)
	}

	headers := r.Group("/headers")
	{
		headers.StaticFile("/file", "./static/headers.dat")
	}

	bch := r.Group("/bch")
	{
		bch.POST("/assemble", api.AssembleBCHTx)
		bch.POST("/broadcast", api.BroadcastBCHTx)
		bch.GET("/hists", api.BchHistory)
	}

	listen := r.Group("/ws")
	{
		listen.GET("/balance", api.NotifyBalanceUpdated)
	}

	version := r.Group("/device")
	{
		version.GET("/update/:device", api.UpdateSoftware)
	}

	notify := r.Group("/notify")
	{
		notify.POST("/", api.GetNotification)
	}

	r.StaticFile("/static", "./bone")

	return r
}
