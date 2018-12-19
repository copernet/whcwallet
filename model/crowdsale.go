package model

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/copernet/whccommon/model"
	"github.com/shopspring/decimal"
)

type CrowdSaleList struct {
	PropertyData string `json:"property_data"`
	BlockTime    int64  `json:"block_time"`
}

func GetAllActiveCrowdSale(keyword string, pagesize int, pageno int) ([]CrowdSaleList, error) {
	var rows *sql.Rows
	var err error
	if keyword != "" {
		pid, err := strconv.Atoi(keyword)
		if err != nil {
			rows, err = db.Table("smart_properties").
				Select("smart_properties.property_data, txes.block_time").
				Where("smart_properties.property_data -> '$.active' = true AND property_data -> '$.name' LIKE ?", "%"+keyword+"%").
				Joins("INNER JOIN txes ON txes.tx_id = smart_properties.create_tx_id").
				Limit(pagesize).
				Offset(pageno*pagesize - pagesize).
				Order("txes.tx_block_height, txes.block_time DESC").
				Rows()
		} else {
			rows, err = db.Table("smart_properties").
				Select("smart_properties.property_data, txes.block_time").
				Where("smart_properties.property_data -> '$.active' = true AND (property_data -> '$.name' LIKE ? OR property_data -> '$.propertyid' = ?)", "%"+keyword+"%", pid).
				Joins("INNER JOIN txes ON txes.tx_id = smart_properties.create_tx_id").
				Limit(pagesize).
				Offset(pageno*pagesize - pagesize).
				Order("txes.tx_block_height, txes.block_time DESC").
				Rows()
		}
	} else {
		rows, err = db.Table("smart_properties").
			Select("smart_properties.property_data, txes.block_time").
			Where("smart_properties.property_data -> '$.active' = true").
			Joins("INNER JOIN txes ON txes.tx_id = smart_properties.create_tx_id").
			Limit(pagesize).
			Offset(pageno*pagesize - pagesize).
			Order("txes.tx_block_height, txes.block_time DESC").
			Rows()
	}

	if err != nil {
		return nil, err
	}

	ret := make([]CrowdSaleList, 0)
	for rows.Next() {
		var item string
		var blockTime int64
		err = rows.Scan(&item, &blockTime)
		if err != nil {
			return nil, err
		}

		ret = append(ret, CrowdSaleList{
			PropertyData: item,
			BlockTime:    blockTime,
		})
	}

	return ret, nil
}

func GetAllActiveCrowdSaleCount(keyword string) (int, error) {
	var total int
	var err error
	if keyword != "" {
		pid, err := strconv.Atoi(keyword)
		if err != nil {
			err = db.Table("smart_properties").
				Where("property_data -> '$.active' = true AND property_data -> '$.name' LIKE ?", "%"+keyword+"%").
				Count(&total).Error
		} else {
			err = db.Table("smart_properties").
				Where("property_data -> '$.active' = true AND (property_data -> '$.name' LIKE ? OR property_data -> '$.propertyid' = ?)", "%"+keyword+"%", pid).
				Count(&total).Error
		}

	} else {
		err = db.Table("smart_properties").
			Where("property_data -> '$.active' = true").
			Count(&total).Error
	}

	if err != nil {
		return 0, err
	}

	return total, nil
}

type PurchaseCrowdSaleRecode struct {
	Fee                        string      `json:"fee"`
	Txid                       string      `json:"txid"`
	Type                       string      `json:"type"`
	Block                      int         `json:"block"`
	Valid                      bool        `json:"valid"`
	Amount                     string      `json:"amount"`
	Ismine                     bool        `json:"ismine"`
	Version                    int         `json:"version"`
	Subsends                   interface{} `json:"subsends"`
	TypeInt                    int         `json:"type_int"`
	Blockhash                  string      `json:"blockhash"`
	Blocktime                  int         `json:"blocktime"`
	Precision                  string      `json:"precision"`
	Propertyid                 int         `json:"propertyid"`
	Issuertokens               string      `json:"issuertokens"`
	Confirmations              int         `json:"confirmations"`
	Actualinvested             string      `json:"actualinvested"`
	Sendingaddress             string      `json:"sendingaddress"`
	Purchasedtokens            string      `json:"purchasedtokens"`
	Referenceaddress           string      `json:"referenceaddress"`
	Purchasedpropertyid        int         `json:"purchasedpropertyid"`
	Purchasedpropertyname      string      `json:"purchasedpropertyname"`
	Purchasedpropertyprecision string      `json:"purchasedpropertyprecision"`
}

func ListPurchaseCrowdsaleTxes(pid uint64, pageSize, pageNo int) ([]PurchaseCrowdSaleRecode, error) {
	rows, err := db.Raw("SELECT txj.tx_data FROM addresses_in_txes as adx "+
		"INNER JOIN tx_jsons as txj "+
		"ON adx.tx_id = txj.tx_id "+
		"WHERE property_id = ? AND address_role = ?",
		pid, model.Participant).
		Limit(pageSize).
		Offset(pageSize*pageNo - pageSize).
		Rows()

	if err != nil {
		return nil, err
	}

	var list []PurchaseCrowdSaleRecode
	var item PurchaseCrowdSaleRecode
	var tmpReceiver string
	for rows.Next() {
		err := rows.Scan(&tmpReceiver)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(tmpReceiver), &item)
		if err != nil {
			return nil, err
		}

		// remove the number string with `0` appended
		amount, err := decimal.NewFromString(item.Amount)
		if err != nil {
			return nil, err
		}
		item.Amount = amount.String()

		issuertokens, err := decimal.NewFromString(item.Issuertokens)
		if err != nil {
			return nil, err
		}
		item.Issuertokens = issuertokens.String()

		actualInvested, err := decimal.NewFromString(item.Actualinvested)
		if err != nil {
			return nil, err
		}
		item.Actualinvested = actualInvested.String()

		purchasedCrowdSale, err := decimal.NewFromString(item.Purchasedtokens)
		if err != nil {
			return nil, err
		}
		item.Purchasedtokens = purchasedCrowdSale.String()

		list = append(list, item)
	}

	return list, nil
}

func GetPurchasedCrowdsaleNumber(pid uint64) (int, error) {
	var counts int
	err := db.Table("addresses_in_txes").Where("property_id = ? AND address_role = ?", pid, model.Participant).Count(&counts).Error
	if err != nil {
		return 0, err
	}

	return counts, nil
}
