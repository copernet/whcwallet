package model

import (
	"time"

	"github.com/copernet/whccommon/model"
	"github.com/shopspring/decimal"
)

func StoreNotification(addr string, item *model.AddressesInTx) error {
	var notification model.Notification
	notification.Address = addr
	notification.Timestamp = time.Now().Unix()
	notification.HistoryID = item.ID

	err := db.Table("notifications").Save(&notification).Error

	return err
}

type APINotify struct {
	Address                     string            `json:"address"`
	PropertyID                  int64             `json:"property_id"`
	PropertyName                string            `json:"property_name"`
	Protocol                    model.Protocol    `json:"protocol"`
	AddressRole                 model.AddressRole `json:"address_role"`
	BalanceAvailableCreditDebit *decimal.Decimal  `json:"balance_available_credit_debit"`
	BalanceReservedCreditDebit  *decimal.Decimal  `json:"balance_reserved_credit_debit"`
	BalanceAcceptedCreditDebit  *decimal.Decimal  `json:"balance_accepted_credit_debit"`
	BalanceFrozenCreditDebit    *decimal.Decimal  `json:"balance_frozen_credit_debit"`
	Timestamp                   int64             `json:"block_time"`
	TxTypeName                  string            `json:"tx_type_name"`
}

type AddressesInTxInherit struct {
	TxTypeName      string
	PropertyName    string
	BlockTime  		int64
	model.AddressesInTx
}


func FilterNotification(addresses []string, from int64, to int64) ([]APINotify, error) {
	if to == 0 {
		to = time.Now().Unix()
	}

	rows, err := db.Table("addresses_in_txes").
		Select("addresses_in_txes.*, txj.tx_data ->> '$.type' as tx_type_name, property_name, txes.block_time").
		Joins("INNER JOIN notifications ON notifications.history_id = addresses_in_txes.id INNER JOIN tx_jsons as txj ON txj.tx_id = addresses_in_txes.tx_id  " +
			"INNER JOIN smart_properties as sp ON sp.property_id = addresses_in_txes.property_id " +
			"AND notifications.address IN (?) " +
			"INNER JOIN txes ON txes.tx_id = txj.tx_id AND txes.block_time BETWEEN ? and ? order by addresses_in_txes.id desc", addresses, from, to).
		Rows()
	if err != nil {
		return nil, err
	}

	var ret []APINotify
	for rows.Next() {
		var notify AddressesInTxInherit
		err = db.ScanRows(rows, &notify)
		if err != nil {
			return nil, err
		}

		//pName, err := GetPropertyName(uint64(notify.PropertyID))
		//if err != nil {
		//	return nil, err
		//}

		apiRet := APINotify{
			Address:                     notify.Address,
			PropertyID:                  notify.PropertyID,
			PropertyName:                notify.PropertyName,
			Protocol:                    notify.Protocol,
			AddressRole:                 notify.AddressRole,
			BalanceAvailableCreditDebit: notify.BalanceAvailableCreditDebit,
			BalanceReservedCreditDebit:  notify.BalanceReservedCreditDebit,
			BalanceAcceptedCreditDebit:  notify.BalanceAcceptedCreditDebit,
			BalanceFrozenCreditDebit:    notify.BalanceFrozenCreditDebit,
			Timestamp:                   notify.BlockTime,
			TxTypeName: 				 notify.TxTypeName,

		}
		ret = append(ret, apiRet)
	}

	return ret, nil
}
