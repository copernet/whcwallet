package model

import (
	"github.com/bcext/cashutil"
	"github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/model/view"
)

func getBalanceForAddress(addr string) ([]model.BalanceForAddress, error) {
	// check address format
	_, err := cashutil.DecodeAddress(addr, config.GetChainParam())
	if err != nil {
		return nil, err
	}

	rows, err := db.Raw("select b.property_id, sp.property_name, sp.precision, b.balance_available, b.pendingpos, b.pendingneg from "+
		"(select tmp.property_id, sum(tmp.balance_available) as balance_available, sum(tmp.pendingpos) as pendingpos, sum(tmp.pendingneg) as pendingneg from "+
		"((select property_id, balance_available, 0 as pendingpos, 0 as pendingneg from address_balances where address = ?) "+
		"UNION ALL (select property_id, 0 as balance_available, sum(if(atx.balance_available_credit_debit > 0, atx.balance_available_credit_debit, 0)) as pendingpos, "+
		"sum(if(atx.balance_available_credit_debit < 0, atx.balance_available_credit_debit, 0)) as pendingneg from addresses_in_txes as atx "+
		"INNER JOIN txes ON atx.tx_id = txes.tx_id Where txes.tx_state = ? AND txes.tx_id < ? AND atx.address = ? AND atx.protocol <> ? Group by atx.property_id)) as tmp "+
		"GROUP BY tmp.property_id) b INNER JOIN smart_properties as sp ON b.property_id = sp.property_id AND sp.protocol <> ? ORDER BY b.property_id;",
		addr, model.Pending, -1, addr, model.BitcoinCash, model.Fiat).Rows()

	if err != nil {
		return nil, err
	}

	ret := make([]model.BalanceForAddress, 0)
	for rows.Next() {
		var bal model.BalanceForAddress
		err := db.ScanRows(rows, &bal)
		if err != nil {
			return nil, err
		}

		bal.Address = addr
		ret = append(ret, bal)
	}

	return ret, nil
}

func GetBalanceForAddresses(addrs []string) (map[string][]model.BalanceForAddress, error) {

	balMap := make(map[string][]model.BalanceForAddress)
	for _, addr := range addrs {
		item, err := getBalanceForAddress(addr)
		if err != nil {
			return nil, err
		}

		balMap[addr] = item
	}

	return balMap, nil
}

func ListOwnersCount(propertyId int) (int64, error) {
	var total int64
	var err error
	err = db.Table("address_balances").
		Where("property_id = ?", propertyId).
		Count(&total).Error

	if err != nil {
		return -1, err
	}
	return total, nil

}

func ListOwners(pageSize int, pageNo int, propertyId int) ([]view.AddressBalance, error) {
	rows, err := db.Table("address_balances adr").Select("adr.address,adr.balance_available,if(tx.tx_type=185,1,0) as status").
		Where("adr.property_id = ?", propertyId).
		Joins("JOIN txes tx on adr.last_tx_id = tx.tx_id ").
		Limit(pageSize).
		Offset(pageSize*pageNo - pageSize).
		Rows()

	models := make([]view.AddressBalance, 0)
	for rows.Next() {
		var vo view.AddressBalance
		err = db.ScanRows(rows, &vo)
		if err != nil {
			return nil, err
		}

		models = append(models, vo)
	}

	return models, nil
}
