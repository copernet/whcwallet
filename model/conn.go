package model

import (
	"fmt"
	"os"

	common "github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var walletdb *gorm.DB

func init() {
	var err error
	db, err = common.ConnectDatabase(config.GetConf().DB)
	if err != nil {
		fmt.Printf("initial database error: %v", err)
		os.Exit(1)
	}

	walletdb, err = common.ConnectDatabase(config.GetConf().WalletDB)
	if err != nil {
		fmt.Printf("initial database error: %v", err)
		os.Exit(1)
	}
}
