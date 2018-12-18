package model

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/copernet/whccommon/model"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func QueryTransactions(property_id int64, address string, start, count int64) ([]string, error) {
	var rows *sql.Rows
	var err error
	if address == "" {
		rows, err = db.Raw("select tx_jsons.tx_data as data "+
			"from property_histories ph, tx_jsons "+
			"where ph.tx_id = tx_jsons.tx_id and ph.property_id = ? "+
			"order by ph.id desc "+
			"limit ?,?", property_id, start, count).Rows()
	} else {
		rows, err = db.Raw("select tx_jsons.tx_data as data "+
			"from property_histories ph, tx_jsons "+
			"where ph.tx_id = tx_jsons.tx_id AND ph.property_id = ? AND tx_jsons.tx_data LIKE ? "+
			"order by ph.id desc "+
			"limit ? offset ?", property_id, "%"+address+"%", count, start).Rows()
	}

	if err != nil {
		return nil, err
	}

	dataSet := make([]string, 0)
	for rows.Next() {
		var data string
		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}

		dataSet = append(dataSet, data)
	}

	return dataSet, nil
}

func QueryTotal(property_id int64, address string) (int64, error) {
	var total int64
	var err error
	if address == "" {
		err = db.Table("property_histories").
			Where("property_id = ?", property_id).
			Joins("INNER JOIN tx_jsons ON property_histories.tx_id = tx_jsons.tx_id").
			Count(&total).Error

		if err != nil {
			return -1, err
		}
		return total, nil

	}

	// if address filter is set.
	rows, err := db.Table("property_histories").
		Select("tx_jsons.tx_data").
		Where("property_id = ?", property_id).
		Joins("INNER JOIN tx_jsons ON tx_jsons.tx_id = property_histories.tx_id").
		Rows()
	if err != nil {
		return -1, err
	}

	for rows.Next() {
		var data string
		err := db.ScanRows(rows, &data)
		if err != nil {
			return -1, err
		}

		if strings.Contains(data, address) {
			total++
		}
	}

	return total, nil
}

func GetHistoryDetail(tx_hash string) (string, error) {
	row := db.Table("txes").
		Select("tx_jsons.tx_data").
		Where("txes.tx_hash = ? AND txes.protocol = ?", tx_hash, model.Wormhole).
		Joins("INNER JOIN tx_jsons ON txes.tx_id = tx_jsons.tx_id").
		Row()

	var tx_data string
	err := row.Scan(&tx_data)
	if err != nil {
		return "", err
	}

	return tx_data, nil
}

func GetHistorySpDetail(tx_hash string) (string, error) {
	row := db.Table("txes").
		Select("property_data").
		Joins("INNER JOIN smart_properties ON tx_id = create_tx_id").
		Where("tx_hash = ?", tx_hash).
		Row()

	var propertyData string
	err := row.Scan(&propertyData)
	if err != nil {
		return "", err
	}

	return propertyData, nil
}

type HistoryList struct {
	TxHash                      string          `json:"tx_hash" gorm:"column:tx_hash"`
	TxType                      int32           `json:"tx_type" gorm:"column:tx_type"`
	TxState                     string          `json:"tx_state" gorm:"column:tx_state"`
	Address                     string          `json:"address"`
	AddressRole                 string          `json:"address_role" gorm:"column:address_role"`
	BalanceAvailableCreditDebit decimal.Decimal `json:"balance_available_credit_debit" gorm:"column:balance_available_credit_debit"`
	BlockTime                   int64           `json:"block_time" gorm:"column:block_time"`
	PropertyName                string          `json:"property_name" gorm:"column:property_name"`
	PropertyID                  uint64          `json:"property_id" gorm:"column:property_id"`
}

func GetHistoryList(addrs []string, pid, pageSize, pageNo int) ([]HistoryList, error) {
	var histories []HistoryList
	var rows *sql.Rows
	var err error
	if pid == 0 {
		rows, err = db.Raw("SELECT tx.tx_hash, tx.tx_type, if(tx.tx_type = 68 and atx.balance_available_credit_debit is null,"+
			"'unmature',tx.tx_state) as tx_state, atx.address, atx.address_role,if(tx.tx_type = 68 and atx.balance_available_credit_debit is null,"+
			"atx.balance_frozen_credit_debit,atx.balance_available_credit_debit) as balance_available_credit_debit, tx.block_time, atx.property_id, "+
			"txj.tx_data ->'$.propertyname' as property_name FROM txes as tx  INNER JOIN addresses_in_txes as atx ON tx.tx_id = atx.tx_id "+
			"INNER JOIN tx_jsons as txj ON txj.tx_id = atx.tx_id AND atx.address IN (?) AND txj.protocol = ?  ORDER BY tx.id DESC", addrs, model.Wormhole).
			Limit(pageSize).
			Offset((pageNo - 1) * pageSize).
			Rows()
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = db.Raw("SELECT tx.tx_hash, tx.tx_type, if(tx.tx_type = 68 and atx.balance_available_credit_debit is null,"+
			"'unmature',tx.tx_state) as tx_state, atx.address, atx.address_role, if(tx.tx_type = 68 and atx.balance_available_credit_debit is null,"+
			"atx.balance_frozen_credit_debit,atx.balance_available_credit_debit) as balance_available_credit_debit, tx.block_time, atx.property_id, "+
			"txj.tx_data ->'$.propertyname' as property_name FROM txes as tx INNER JOIN addresses_in_txes as atx ON tx.tx_id = atx.tx_id AND atx.property_id = ? "+
			"INNER JOIN tx_jsons as txj ON txj.tx_id = atx.tx_id  AND atx.address IN (?) AND txj.protocol = ? "+
			"ORDER BY tx.id DESC", pid, addrs, model.Wormhole).
			Limit(pageSize).
			Offset((pageNo - 1) * pageSize).
			Rows()
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var history HistoryList
		err := db.ScanRows(rows, &history)
		if err != nil {
			return nil, err
		}

		history.PropertyName = strings.Trim(history.PropertyName, "\"")

		if history.PropertyName == "" {
			name, err := GetPropertyName(history.PropertyID)
			if err != nil {
				return nil, err
			}

			history.PropertyName = name
		}

		histories = append(histories, history)
	}

	return histories, nil
}

func GetHistoryListCount(addrs []string, pid int) (int, error) {
	var row *sql.Row
	if pid == 0 {
		row = db.Raw("SELECT count(*) as count FROM txes as tx "+
			"INNER JOIN addresses_in_txes as atx ON tx.tx_id = atx.tx_id "+
			"INNER JOIN tx_jsons as txj ON txj.tx_id = atx.tx_id "+
			"AND atx.address IN (?) AND txj.protocol = ? ",
			addrs, model.Wormhole).
			Row()
	} else {
		row = db.Raw("SELECT count(*) as count FROM txes as tx "+
			"INNER JOIN addresses_in_txes as atx ON tx.tx_id = atx.tx_id AND atx.property_id = ? "+
			"INNER JOIN tx_jsons as txj ON txj.tx_id = atx.tx_id "+
			"AND atx.address IN (?) AND txj.protocol = ? ",
			pid, addrs, model.Wormhole).
			Row()
	}

	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func GetNotificationItem(address string, pid int64, txID int) (*model.AddressesInTx, error) {
	var ret model.AddressesInTx
	err := db.Table("addresses_in_txes").
		Select("id").
		Where("address = ? AND tx_id = ? AND property_id = ?", address, txID, pid).
		First(&ret).Error

	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func parsePrecision(prec string) (int, error) {
	reg, err := regexp.Compile("[0-9]+")
	if err != nil {
		return 0, err
	}

	ret := reg.FindStringSubmatch(prec)
	if len(ret) == 0 {
		return 0, errors.New("not found precision")
	}

	precision, err := strconv.Atoi(ret[0])
	if err != nil {
		return 0, err
	}

	return precision, nil
}
