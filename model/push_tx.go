package model

import (
	"time"

	common "github.com/copernet/whccommon/model"
	"github.com/jinzhu/gorm"
)

func Begin() *gorm.DB {
	return db.Begin()
}

func IsExistTx(txid string, protocol common.Protocol) bool {
	var tx common.Tx
	notFound := db.Where("tx_hash = ? AND protocol = ?", txid, protocol).First(&tx).RecordNotFound()

	return !notFound
}

func InsertTx(tx *common.Tx, dbtx *gorm.DB) (int, error) {
	if dbtx == nil {
		action := db.Begin()
		err := action.Exec("INSERT INTO `txes` (`tx_hash`, `tx_id`, `protocol`, `tx_type`, `ecosystem`, `tx_state`, `block_time`) "+
			"select ? as tx_hash, case when min(tx_id) < -1 then min(tx_id)-1 else -2 end as tx_id，? as protocol, ? as tx_type, "+
			"? as ecosystem, ? as tx_state, ? as block_time from txes",
			tx.TxHash, tx.Protocol, tx.TxType, tx.Ecosystem, tx.TxState, time.Now().Unix()).Error
		if err != nil {
			action.Rollback()
			return 0, err
		}

		var txTmp common.Tx
		err = action.Table("txes").Select("tx_id").Where("tx_hash = ?", tx.TxHash).Order("tx_id").First(&txTmp).Error
		if err != nil {
			action.Rollback()
			return 0, err
		}

		return txTmp.TxID, err
	}

	err := dbtx.Exec("INSERT INTO `txes` (`tx_hash`, `tx_id`, `protocol`, `tx_type`, `ecosystem`, `tx_state`, `block_time`) "+
		"select ? as tx_hash, case when min(tx_id) < -1 then min(tx_id)-1 else -2 end as tx_id, ? as protocol, ? as tx_type, "+
		"? as ecosystem, ? as tx_state, ? as block_time from txes",
		tx.TxHash, tx.Protocol, tx.TxType, tx.Ecosystem, tx.TxState, time.Now().Unix()).Error
	if err != nil {
		return 0, err
	}

	var txTmp common.Tx
	err = dbtx.Table("txes").Select("tx_id").Where("tx_hash = ?", tx.TxHash).Order("tx_id").First(&txTmp).Error

	return txTmp.TxID, err
}

func InsertAddressesInTx(addrInTx *common.AddressesInTx, dbtx *gorm.DB) error {
	if dbtx == nil {
		return db.Create(addrInTx).Error
	}

	return dbtx.Create(addrInTx).Error
}

func InsertTxJson(txJson *common.TxJson, dbtx *gorm.DB) error {
	if dbtx == nil {
		return db.Create(txJson).Error
	}

	return dbtx.Create(txJson).Error
}
