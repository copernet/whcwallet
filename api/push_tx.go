package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/copernet/whcwallet/model/view"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/whc.go/btcjson"
	"github.com/copernet/whccommon/log"
	common "github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/logic"
	"github.com/copernet/whcwallet/logic/ws"
	"github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func PushTx(c *gin.Context) {
	signedTx := c.PostForm("signedTx")

	// check post form data required
	if signedTx == "" {
		log.WithCtx(c).Error("broadcast empty encoding transaction")

		c.JSON(200, apiErrorWithMsg(ErrFormItems, "post data lack of signed transaction"))
		return
	}

	rawdata, err := hex.DecodeString(signedTx)
	if err != nil {
		log.WithCtx(c).Errorf("hex decode raw transaction string error: %v", err)

		c.JSON(200, apiError(ErrHexStringFormat))
		return
	}

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(rawdata))
	if err != nil {
		log.WithCtx(c).Errorf("deserialize transaction error: %v", err)

		c.JSON(200, apiError(ErrTxDeserialize))
		return
	}

	// broadcast the post transaction
	client := view.GetRPCIns()
	txhash, err := client.SendRawTransaction(&tx, false)
	if err != nil {
		log.WithCtx(c).Warnf("the transaction posted broadcasts failed: %v", err)

		c.JSON(200, apiErrorWithMsg(ErrSendRawTransaction,
			"Please check carefully and now it is not a correct transaction"))
		return
	}

	// async query firstly.
	detail := client.WhcGetTransactionAsync(txhash.String())

	// Do not get result via rpc request. This function is from btcsuite/btcd.
	decodedTx, err := util.DecodeRawTransaction(signedTx, config.GetChainParam())
	if err != nil {
		log.WithCtx(c).Warnf("the transaction decode error: %v", err)

		c.JSON(200, apiError(ErrDecodeRawTransaction))
		return
	}

	inputs, _, err := getInputs(&tx)
	if err != nil {
		log.WithCtx(c).Warnf("get transactions input error from BCH API: %v", err)

		c.JSON(200, apiError(ErrCanNotGetInputs))
		return
	}

	// start database transaction
	dbtx := model.Begin()

	err = insertBCH(decodedTx, inputs, dbtx)
	if err != nil {
		dbtx.Rollback()

		log.WithCtx(c).Errorf("Insert transaction for relative BCH error: %v", err)

		c.JSON(200, apiError(ErrInsertBCH))
		return
	}

	ret, err := client.WhcDecodeTransaction(signedTx, nil, nil)
	if err != nil {
		log.WithCtx(c).Errorf("whc_decodetransaction rpc request error: %v", err)

		c.JSON(200, apiError(ErrWhcDecodeTransaction))
		return
	}

	detailRec, err := detail.Receive()
	if err != nil {
		log.WithCtx(c).Errorf("whc_decodetransaction rpc request error: %v", err)

		c.JSON(200, apiError(ErrWhcGetTransaction))
		return
	}

	dbtx.Commit()

	err = model.Publish(model.MempoolTxTip, detailRec.TxID)
	log.WithCtx(c).Infof("public tx:%s to topic:%s", detailRec.TxID, model.MempoolTxTip)
	//err = insertWhc(ret, detailRec, dbtx)
	if err != nil {
		log.WithCtx(c).Errorf("Publish transaction error: %v", err)
		c.JSON(200, apiError(ErrInsertWhc))
		return
	}

	// omit wormhole pending balance via websocket
	sender := ret.SendingAddress
	receiver := ret.ReferenceAddress

	senderConnection, ok := connMgr.GetConnForAddress(sender)
	if !ok {
		logrus.Infof("the client holds this address not online: %s", sender)
	} else {
		senderBal, err := GetBalanceFromCache([]string{sender}, true)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				log.DefaultTraceLabel: senderConnection.Uid,
			}).Errorf("get wormhole balance for websocket notification error: %v; address: %s", err, sender)
		} else {
			var msg ws.Message
			msg.Address = sender
			msg.Symbol = ws.CoinWormhole
			msg.Balance = senderBal[sender]
			connMgr.MsgChan <- msg
		}
	}

	if receiver != "" {
		receiverConnection, ok := connMgr.GetConnForAddress(receiver)
		if !ok {
			logrus.Infof("the client holds this address not online: %s", receiver)
		} else {
			receiverBal, err := GetBalanceFromCache([]string{receiver}, true)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					log.DefaultTraceLabel: receiverConnection.Uid,
				}).Errorf("get wormhole balance for websocket notification error: %v; address: %s", err, receiver)
			} else {
				var msg ws.Message
				msg.Address = receiver
				msg.Symbol = ws.CoinWormhole
				msg.Balance = receiverBal[receiver]
				connMgr.MsgChan <- msg
			}
		}
	}

	c.JSON(200, apiSuccess(map[string]string{
		"txhash": txhash.String(),
	}))
}

func getInputs(tx *wire.MsgTx) (map[string]int64, bool, error) {
	node := getNode()
	var txid string
	var vout uint32
	var valid bool
	inputs := make(map[string]int64)
	for _, txin := range tx.TxIn {
		txid = txin.PreviousOutPoint.Hash.String()
		vout = txin.PreviousOutPoint.Index
		detail, err := node.BlockchainTransactionGet(txid, true)
		if err != nil {
			return nil, false, err
		}

		out := detail.Vout[vout]
		value, err := cashutil.NewAmount(out.Value)
		if err != nil {
			return nil, false, err
		}

		if out.ScriptPubKey.Type != "P2PKH" && out.ScriptPubKey.Type != "P2SH" {
			valid = false
		}

		for _, item := range detail.Vout[vout].ScriptPubKey.Addresses {
			//The address item from api result is base58 encoded format, and we
			// should convert to bech32 encoded format for the following using.
			addr, _ := cashutil.DecodeAddress(item, config.GetChainParam())
			addrStr := addr.EncodeAddress(true)
			if amount, ok := inputs[addrStr]; ok {
				inputs[addrStr] = amount + int64(value)
			} else {
				inputs[addrStr] = int64(value)
			}

		}
	}

	return inputs, valid, nil
}

func insertBCH(decodedTx *btcjson.TxRawDecodeResult, inputs map[string]int64, dbtx *gorm.DB) error {
	var pid int64
	txtype := 0
	protocol := common.BitcoinCash

	// if tx have existed in database, skip
	exist := model.IsExistTx(decodedTx.Txid, protocol)
	if exist {
		return errors.New("the bitcoin cash transaction has existed")
	}

	var addressTxIndex int

	pendingTx := common.Tx{
		TxHash:    decodedTx.Txid,
		Protocol:  protocol,
		TxType:    uint64(txtype),
		Ecosystem: common.Production,
		TxState:   common.Pending,
	}
	minTxID, err := model.InsertTx(&pendingTx, dbtx)
	if err != nil {
		return err
	}

	for address, amount := range inputs {
		value := cashutil.Amount(amount)

		change := decimal.New(int64(value), 0).Div(decimal.New(1e8, 0))
		addressInTx := common.AddressesInTx{
			Address:                     address,
			PropertyID:                  pid,
			Protocol:                    protocol,
			TxID:                        minTxID,
			AddressTxIndex:              int16(addressTxIndex),
			AddressRole:                 common.Sender,
			BalanceAvailableCreditDebit: &change,
		}
		if err := model.InsertAddressesInTx(&addressInTx, dbtx); err != nil {
			return err
		}

		addressTxIndex++
	}

	addressTxIndex = 0
	for _, out := range decodedTx.Vout {
		// Here, we do not need to deal with the error, because the tx has been
		// deserialized correctly. That is to say, the tx's amount is correct.
		value, _ := cashutil.NewAmount(out.Value)
		change := decimal.New(int64(value), 0).Div(decimal.New(1e8, 0))
		if out.ScriptPubKey.Type != "nulldata" {
			for _, addr := range out.ScriptPubKey.Addresses {
				addrTxIndex := common.AddressesInTx{
					Address:                     addr,
					PropertyID:                  pid,
					Protocol:                    protocol,
					TxID:                        minTxID,
					AddressTxIndex:              int16(addressTxIndex),
					AddressRole:                 common.Recipient,
					BalanceAvailableCreditDebit: &change,
				}

				if err := model.InsertAddressesInTx(&addrTxIndex, dbtx); err != nil {
					return err
				}
			}

			addressTxIndex++
		}
	}

	ret, err := json.Marshal(decodedTx)
	if err != nil {
		return err
	}
	// store signed tx until it confirms
	txjson := common.TxJson{
		TxID:     minTxID,
		Protocol: protocol,
		TxData:   string(ret),
	}
	if err := model.InsertTxJson(&txjson, dbtx); err != nil {
		return err
	}

	return nil
}

func insertWhc(whcDecodeTx *btcjson.GenerateTransactionResult, whcGetTx *btcjson.GenerateTransactionResult, dbtx *gorm.DB) error {
	sender := whcDecodeTx.SendingAddress
	receiver := whcDecodeTx.ReferenceAddress

	var pid int64
	if whcDecodeTx.PropertyID != 0 {
		pid = whcDecodeTx.PropertyID
	}

	txid := whcDecodeTx.TxID
	// if tx is already in database, skip
	if ok := model.IsExistTx(txid, common.Wormhole); ok {
		return errors.New("the wormhole transaction has existed")
	}

	var amount cashutil.Amount
	if whcDecodeTx.Amount != "" {
		numDecimal, err := decimal.NewFromString(whcDecodeTx.Amount)
		if err != nil {
			return err
		}

		numFloat, _ := numDecimal.Float64()
		amount, _ = cashutil.NewAmount(numFloat)
	}

	var sendAmount cashutil.Amount
	var recvAmount cashutil.Amount

	switch whcDecodeTx.TypeInt {
	case 55, 68:
		sendAmount = 0
		recvAmount = amount
	case 50, 51, 54:
		sendAmount = amount
		recvAmount = 0
	default:
		sendAmount = -amount
		recvAmount = amount
	}

	tx := common.Tx{
		TxHash:    txid,
		Protocol:  common.Wormhole,
		TxType:    uint64(whcDecodeTx.TypeInt),
		Ecosystem: common.Production,
		TxState:   common.Pending,
	}
	minTxID, err := model.InsertTx(&tx, dbtx)
	if err != nil {
		return err
	}

	// insert the addressesintxs entry for the sender
	if sendAmount != 0 || whcDecodeTx.TypeInt == 70 {
		change := decimal.New(int64(sendAmount), 0).Div(decimal.New(1e8, 0))
		addressInTxs := common.AddressesInTx{
			Address:                     sender,
			PropertyID:                  pid,
			Protocol:                    common.Wormhole,
			TxID:                        minTxID,
			AddressTxIndex:              0,
			BalanceAvailableCreditDebit: &change,

			// the following fields uses default value
			AddressRole: common.Sender,
			// BalanceAcceptedCreditDebit:0
		}

		switch whcDecodeTx.TypeInt {
		case 70:
			// for change issuer transaction record
			// todo validate the role field
			addressInTxs.AddressRole = common.Seller
		case 50, 51, 54:
			addressInTxs.AddressRole = common.Issuer
		}

		if err := model.InsertAddressesInTx(&addressInTxs, dbtx); err != nil {
			return err
		}
	}

	if receiver != "" && recvAmount != 0 {
		change := decimal.New(int64(recvAmount), 0).Div(decimal.New(1e8, 0))
		addrInTxs := common.AddressesInTx{
			Address:                     receiver,
			PropertyID:                  pid,
			Protocol:                    common.Wormhole,
			TxID:                        minTxID,
			AddressTxIndex:              0,
			BalanceAvailableCreditDebit: &change,

			// the following fields uses default value
			AddressRole: common.Recipient,
			// BalanceAcceptedCreditDebit:0,
		}

		if err := model.InsertAddressesInTx(&addrInTxs, dbtx); err != nil {
			return err
		}
	}

	ret, err := json.Marshal(whcGetTx)
	if err != nil {
		return err
	}
	// store signed tx until it confirms
	txjson := common.TxJson{
		TxID:     minTxID,
		Protocol: common.Wormhole,
		TxData:   string(ret),
	}
	if err := model.InsertTxJson(&txjson, dbtx); err != nil {
		return err
	}

	return nil
}

func FeeRate(c *gin.Context) {
	//Load cache feeRate
	rates, err := model.GetFeeRate()
	if err != nil {
		if strings.Contains(err.Error(), "nil returned") {
			log.WithCtx(c).Info("the stored feerate in redis has timeout")
		} else {
			log.WithCtx(c).Errorf("fetch feerate from redis error:%s", err.Error())
		}

		rates = logic.EstimateFee(c)
		model.StoreFeeRate(rates)
	}

	c.JSON(200, apiSuccess(rates))
}
