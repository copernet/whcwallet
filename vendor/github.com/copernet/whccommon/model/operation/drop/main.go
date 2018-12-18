package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/copernet/whccommon/model"
)

func main() {
	host := flag.String("host", "127.0.0.1", "please input database connection host")
	port := flag.Int("port", 3306, "please input database connection port")
	user := flag.String("user", "root", "please input database administrator username")
	passwd := flag.String("passwd", "", "please input database administrator password")
	db := flag.String("db", "wormhole", "please input the common database name['-' represents ignoring this database]")
	flag.Parse()

	config := model.DBOption{
		Host:     *host,
		Port:     *port,
		User:     *user,
		Passwd:   *passwd,
		Database: *db,
		Log:      true,
	}

	if err := dropDatabase(&config); err != nil {
		fmt.Printf("Unfortunately: %v", err)
		os.Exit(1)
	}
}

func dropDatabase(config *model.DBOption) error {
	db, err := model.ConnectDatabase(config)
	if err != nil {
		return err
	}
	defer db.Close()

	db.DropTableIfExists(
		&model.Block{},
		&model.Tx{},
		&model.TxJson{},
		&model.AddressesInTx{},
		&model.AddressBalance{},
		&model.SmartProperty{},
		&model.PropertyHistory{},
		&model.Session{})

	return nil
}
