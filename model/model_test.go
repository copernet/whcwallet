package model

import (
	"testing"

	"github.com/copernet/whccommon/model"
)

func TestScanFunc(t *testing.T) {
	rows, err := db.Table("txes").Select([]string{"id", "tx_type"}).Rows()
	if err != nil {
		t.Error(err)
	}

	txs := make([]model.Tx, 0)

	// define receiver variable at here, do not influence the result set,
	// that is to say that the result set will not be one record.
	var id uint
	var txtype uint64

	for rows.Next() {
		err = rows.Scan(&id, &txtype)
		if err != nil {
			t.Error(err)
		} else {
			txs = append(txs, model.Tx{
				ID:     id, // because a // ssignment via the value
				TxType: txtype,
			})
		}

	}
}
