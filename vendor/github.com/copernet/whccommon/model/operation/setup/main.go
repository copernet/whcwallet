package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/copernet/whccommon/model"
	"github.com/jinzhu/gorm"
)

func main() {
	host := flag.String("host", "127.0.0.1", "please input database connection host")
	port := flag.Int("port", 3306, "please input database connection port")
	user := flag.String("user", "root", "please input database administrator username")
	passwd := flag.String("passwd", "", "please input database administrator password")
	db := flag.String("db", "wormhole", "please input the common database name('-' represents ignoring this database)")
	walletdb := flag.String("walletdb", "walletdb", "please input the walletdb database name('-' represents ignoring this database)")
	flag.Parse()

	config := model.DBOption{
		Host:   *host,
		Port:   *port,
		User:   *user,
		Passwd: *passwd,
		Log:    true,
	}

	if *db != "-" {
		config.Database = *db

		db, err := model.ConnectDatabase(&config)
		if err != nil {
			fmt.Printf("Unfortunately: %v", err)
			os.Exit(1)
		}
		defer db.Close()

		createTable(db, &model.Block{})
		createTable(db, &model.Tx{})
		createTable(db, &model.TxJson{})
		createTable(db, &model.AddressesInTx{})
		createTable(db, &model.AddressBalance{})
		createTable(db, &model.SmartProperty{})
		createTable(db, &model.PropertyHistory{})
		createTable(db, &model.Version{})
		createTable(db, &model.Notification{})
	}

	if *walletdb != "-" {
		config.Database = *walletdb

		db, err := model.ConnectDatabase(&config)
		if err != nil {
			fmt.Printf("Unfortunately: %v", err)
			os.Exit(1)
		}
		defer db.Close()

		createTable(db, &model.Session{})
		createTable(db, &model.Wallet{})
	}

}

func createTable(db *gorm.DB, table interface{}) {
	// Migrate the schema
	db.AutoMigrate(table)

	// Create
	if !db.HasTable(table) {
		db.CreateTable(table)
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(table)
	}
}
