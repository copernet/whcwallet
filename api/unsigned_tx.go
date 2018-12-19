package api

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/whc.go/btcjson"
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/logic"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const (
	defaultInputSize     = 148
	defaultOutputSize    = 34
	defaultSignatureSize = 107
	defaultFeeRate = 0.00001148
)

var bchUnit = decimal.New(1e8, 0)

func GetUnsignedTx(c *gin.Context) {
	txStr := c.Param("txtype")
	txType, err := strconv.Atoi(txStr)
	if err != nil {
		c.JSON(200, apiError(ErrConvertInt))
		return
	}

	var unsignedTx string
	var signData []btcjson.PrevTx

	txBuilder := logic.FactoryForTxCheckAndCreate(txType)
	err = txBuilder.Check(c)
	if err != nil {
		log.WithCtx(c).Errorf("user input validation failed: %v", err)

		// judge whether the error is exceed nulldata length limit or not
		if e, ok := err.(*logic.ErrNulldataLength); ok {
			c.JSON(200, &Response{
				Code:    ErrExceedMaxNulldataLimit,
				Message: e.Msg,
			})

			return
		}

		log.WithCtx(c).Errorf("user post invalid form: %v", err)

		c.JSON(200, apiErrorWithMsg(ErrFormItems, err.Error()))
		return
	}

	payload, err := txBuilder.CreatePayload(c)
	if err != nil {
		log.WithCtx(c).Errorf("create payload error: %v", err)

		c.JSON(200, apiError(ErrCreatePayload))
		return
	}

	unsignedTx, signData, err = assembleTx(c, payload)
	if err != nil {
		log.WithCtx(c).Errorf("assemble transaction failed: %v", err)
		return
	}

	c.JSON(200, apiSuccess(map[string]interface{}{
		"unsigned_tx": unsignedTx,
		"sign_data":   signData,
	}))
	return
}

var emptyError = errors.New("empty error")

func assembleTx(c *gin.Context, payload string) (string, []btcjson.PrevTx, error) {
	precision := c.PostForm("precision")
	if precision != "" && !IsSupportPrecision(precision) {
		c.JSON(200, apiError(ErrUnSupportPrecision))
		return "", nil, emptyError
	}

	feeRate, err := decimal.NewFromString(c.PostForm("fee"))
	if err != nil {
		c.JSON(200, apiError(ErrIncorrectAmount))
		return "", nil, err
	}


	// to keep max 8 digital number for feerate
	feeRate = feeRate.Truncate(8)
	fee := decimal.NewFromFloat(defaultFeeRate)

	if feeRate.LessThan(fee) {
		feeRate = fee
	}

	// txtype must be existed, checked at upstream
	txType, _ := c.Get("txtype")

	var totalOutput decimal.Decimal
	var burnAmount decimal.Decimal

	// transaction for creating whc
	if txType == 68 {
		burnAmount, err = decimal.NewFromString(c.PostForm("amount_for_burn"))
		if err != nil {
			c.JSON(200, apiError(ErrIncorrectAmount))
			return "", nil, emptyError
		}
		totalOutput = totalOutput.Add(burnAmount)
	}

	to := c.PostForm("transaction_to")
	var toAddr cashutil.Address
	if to != "" {
		toAddr, err = cashutil.DecodeAddress(to, config.GetChainParam())
		if err != nil {
			c.JSON(200, apiError(ErrIncorrectAddress))
			return "", nil, err
		}
		refAmountStr := c.DefaultPostForm("reference_amount", config.GetConf().Tx.MiniOutput)
		refAmount, err := decimal.NewFromString(refAmountStr)
		if err != nil {
			c.JSON(200, apiError(ErrIncorrectAmount))
			return "", nil, err
		}

		totalOutput = totalOutput.Add(refAmount)
	}

	from := c.PostForm("transaction_from")
	addr, err := cashutil.DecodeAddress(from, config.GetChainParam())
	if err != nil {
		c.JSON(200, apiError(ErrIncorrectAddress))
		return "", nil, err
	}

	// rpc instance
	client := view.GetRPCIns()

	var rawTx string
	if txType == 68 {
		burnAmountStr := burnAmount.String()
		rawTx, err = client.WhcCreateRawTxReference(rawTx, config.GetBurningAddress(), &burnAmountStr)
		if err != nil {
			c.JSON(200, apiError(ErrCreateRawTxReference))
			return "", nil, err
		}

		rawTx, err = client.WhcCreateRawTxOpReturn(rawTx, payload)
		if err != nil {
			c.JSON(200, apiError(ErrCreateRawTxOpReturn))
			return "", nil, err
		}
	} else {
		if to != "" {
			refAmountStr := c.DefaultPostForm("reference_amount", config.GetConf().Tx.MiniOutput)
			rawTx, err = client.WhcCreateRawTxReference(rawTx, toAddr.EncodeAddress(true), &refAmountStr)
			if err != nil {
				c.JSON(200, apiError(ErrCreateRawTxReference))
				return "", nil, err
			}
		}

		rawTx, err = client.WhcCreateRawTxOpReturn(rawTx, payload)
		if err != nil {
			c.JSON(200, apiError(ErrCreateRawTxOpReturn))
			return "", nil, err
		}
	}

	utxoList, err := GetUtxoElectrumx(addr, decimal.Zero, true)
	if err != nil {
		c.JSON(200, apiErrorWithMsg(ErrCanNotGetUtxo, "get spendable coins error: "+err.Error()))
		return "", nil, err
	}

	// now the partial transaction serialize size is known:
	// total: len(rawTx)/2 + 34 bytes[change output],
	// add one input => add 148 bytes for transaction.
	baseSize := len(rawTx) / 2
	baseSizeWithChange := len(rawTx)/2 + defaultOutputSize
	var totalInput float64
	var isNoChange, abundant bool
	var finalUtxoList []btcjson.PrevTx

	// attempt to create a transaction without change output
	for _, utxo := range utxoList {
		baseSize += defaultInputSize
		totalInput += utxo.Value
		fee := feeRate.Mul(decimal.NewFromFloat(float64(baseSize) / 1000))
		diff := decimal.NewFromFloat(totalInput).Sub(fee.Add(totalOutput)).Mul(bchUnit).IntPart()
		finalUtxoList = append(finalUtxoList, utxo)
		// the diff serve fee
		if diff >= 0 && diff < dustSatoshi {
			abundant = true
			isNoChange = true
			goto next
			break
		}
	}

	// reset value
	totalInput = 0
	isNoChange = false
	abundant = false
	finalUtxoList = []btcjson.PrevTx{}

	// there should have a change output
	for _, utxo := range utxoList {
		baseSizeWithChange += defaultInputSize
		totalInput += utxo.Value
		fee := feeRate.Mul(decimal.NewFromFloat(float64(baseSizeWithChange) / 1000))
		diff := decimal.NewFromFloat(totalInput).Sub(fee.Add(totalOutput)).Mul(bchUnit).IntPart()
		finalUtxoList = append(finalUtxoList, utxo)
		if diff >= dustSatoshi {
			abundant = true
			break
		}
	}

next:

	if !abundant {
		c.JSON(200, apiError(ErrInsufficientBalance))
		return "", nil, errors.New("Account has insufficient balance fro creating transaction")
	}

	// insert all utxo list
	// TODO not via rpc, assemble by myself
	for _, utxo := range finalUtxoList {
		rawTx, err = client.WhcCreateRawTxInput(rawTx, utxo.TxID, int(utxo.Vout))
		if err != nil {
			c.JSON(200, apiError(ErrCreateRawTxInput))
			return "", nil, err
		}
	}

	var tx wire.MsgTx
	b, _ := hex.DecodeString(rawTx)
	err = tx.Deserialize(bytes.NewBuffer(b))
	if err != nil {
		c.JSON(200, apiError(ErrTxDeserialize))
		return "", nil, err
	}

	if !isNoChange {
		finalFee := feeRate.Mul(decimal.New(int64(len(rawTx)/2+len(finalUtxoList)*defaultSignatureSize+defaultOutputSize), -3))
		realOutput, _ := totalOutput.Add(finalFee).Float64()
		changeAmount := int64((totalInput - realOutput) * 1e8)

		if txType == 68 {
			redeemAddress := c.DefaultPostForm("redeem_address", addr.EncodeAddress(true))
			address, err := cashutil.DecodeAddress(redeemAddress, config.GetChainParam())
			if err != nil {
				c.JSON(200, apiError(ErrIncorrectAddress))
				return "", nil, err
			}

			pkScript, _ := txscript.PayToAddrScript(address)
			changeOut := wire.NewTxOut(changeAmount, pkScript)
			tx.TxOut = append(tx.TxOut, changeOut)
		} else {
			redeemAddress := c.DefaultPostForm("redeem_address", from)
			redeem, err := cashutil.DecodeAddress(redeemAddress, config.GetChainParam())
			if err != nil {
				c.JSON(200, apiError(ErrIncorrectAddress))
				return "", nil, err
			}

			pkScript, _ := txscript.PayToAddrScript(redeem)
			changeOut := wire.NewTxOut(changeAmount, pkScript)
			tx.TxOut = append(tx.TxOut, changeOut)
			tx.TxOut[0], tx.TxOut[len(tx.TxOut)-1] = tx.TxOut[len(tx.TxOut)-1], tx.TxOut[0]
		}

		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		tx.Serialize(buf)
		rawTx = hex.EncodeToString(buf.Bytes())
	}

	return rawTx, finalUtxoList, nil
}

func GetUtxoElectrumx(addr cashutil.Address, requireAmount decimal.Decimal, needAllUtxos bool) ([]btcjson.PrevTx, error) {
	// check the balance if sufficient
	addrStr := addr.EncodeAddress(false)
	balance, err := getNode().BlockchainAddressGetBalance(addrStr)
	if err != nil {
		return nil, err
	}

	exp := decimal.New(1e8, 0)
	var requiredInt int64 = 0
	if !needAllUtxos {
		requiredInt = requireAmount.Mul(exp).IntPart()
	}

	balanceInt := int64(balance.Confirmed)
	if balanceInt < requiredInt {
		str := "insufficient balance from address: " + addrStr
		return nil, errors.New(str)
	}

	// find the valid utxos
	utxos, err := getNode().BlockchainAddressListUnspent(addrStr)
	ret := make([]btcjson.PrevTx, 0, defaultUnspentList)
	if utxos != nil && len(utxos) > 0 && err == nil {
		for _, utxo := range utxos {
			if requiredInt > 0 || needAllUtxos {
				tx, err := getNode().BlockchainTransactionGet(utxo.Hash, true)
				if err == nil {
					// coinbase utxo must have 100 confirmations
					isCoinBase := false
					if tx.Vin[0].Coinbase != "" {
						isCoinBase = true
					}

					confirmations := tx.Confirmations
					if isCoinBase && confirmations < 100 {
						continue
					}
					// only count the pubkeyhash utxo, for p2sh utxo, we can't check if it's a multisig utxo while it hasn't been spent
					scriptType := tx.Vout[utxo.Pos].ScriptPubKey.Type
					if scriptType != "pubkeyhash" {
						continue
					}
					hexScript := tx.Vout[utxo.Pos].ScriptPubKey.Hex
					requiredInt -= int64(utxo.Value)
					//construct the utxo obj to return
					amount := cashutil.Amount(utxo.Value)
					coin := btcjson.PrevTx{
						TxID:         utxo.Hash,
						Vout:         int(utxo.Pos),
						Value:        amount.ToBCH(),
						ScriptPubKey: hexScript,
					}
					ret = append(ret, coin)
				}
			} else {
				break
			}
		}
		if requiredInt > 0 && !needAllUtxos {
			str := "insufficient balance from address: " + addrStr + ", balance is sufficient, but valid utxos are not enough"
			return nil, errors.New(str)
		} else {
			return ret, nil
		}
	} else {
		return nil, err
	}
}

func getUtxo(addr cashutil.Address, requireAmount decimal.Decimal, page int,
	ret []btcjson.PrevTx, avail int64) ([]btcjson.PrevTx, error) {

	addrStr := addr.EncodeAddress(false)
	balance, err := getBalance(addrStr)
	if err != nil {
		return nil, err
	}

	exp := decimal.New(1e8, 0)
	if balance < requireAmount.Mul(exp).IntPart() {
		str := "insufficient balance from address: " + addrStr
		return nil, errors.New(str)
	}

	content, err := getUnspent(addrStr, page)
	if err != nil {
		return nil, err
	}

	// encounter an error
	if gjson.Get(content, "err_no").String() != "0" {
		return nil, errors.New(gjson.Get(content, "err_msg").String())
	}

	lists := gjson.Get(content, "data.list").Array()
	// not unspent coins
	if len(lists) == 0 {
		return nil, errors.New("not available balance")
	}

	// initial ret container
	if ret == nil {
		ret = make([]btcjson.PrevTx, 0, defaultUnspentList)
	}

	for _, item := range lists {
		// stop iterate if have sufficient balance
		if avail > int64(requireAmount.Mul(exp).IntPart()) {
			break
		}

		hash := item.Get("tx_hash").String()
		index := item.Get("tx_output_n").Int()
		// check whether the specified transaction output is spendable or not
		ok, script, err := isAvailableCoin(hash, int(index))
		if !ok || err != nil {
			continue
		}

		value := item.Get("value").Int()
		amount := cashutil.Amount(value)
		coin := btcjson.PrevTx{
			TxID:         hash,
			Vout:         int(index),
			Value:        amount.ToBCH(),
			ScriptPubKey: script,
		}
		ret = append(ret, coin)

		// update total avail
		avail += value
	}

	if avail < int64(requireAmount.Mul(exp).IntPart()) &&
		gjson.Get(content, "total_count").Int()-defaultPageSize*int64(page) > 0 {
		// iterate the self function
		return getUtxo(addr, requireAmount, page+1, ret, avail)
	} else {
		return ret, nil
	}
}

// get balance for the specified address
func getBalance(addr string) (int64, error) {
	url := config.GetBCHAPI() + "/address/" + addr
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	return gjson.Get(string(content), "data.balance").Int(), nil
}

// get raw string of unspent list for the specified address
func getUnspent(addr string, page int) (string, error) {
	url := config.GetBCHAPI() + "/address/" + addr + "/unspent?pagesize=" +
		strconv.Itoa(defaultPageSize) + "&page=" + strconv.Itoa(page)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func isAvailableCoin(txid string, index int) (bool, string, error) {
	url := config.GetBCHAPI() + "/tx/" + txid + "?verbose=3"
	res, err := http.Get(url)
	if err != nil {
		return false, "", err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, "", err
	}

	ret := string(content)

	// check error from API
	if gjson.Get(ret, "err_no").String() != "0" {
		str := "get transaction error from BCH API"
		return false, "", errors.New(str)
	}

	// check coinbase mature
	if gjson.Get(ret, "data.is_coinbase").Bool() &&
		gjson.Get(ret, "date.confirmations").Int() < 100 {
		return false, "", errors.New("can not spent immature coinbase transaction")
	}

	// check the transaction type of the specified output
	if gjson.Get(ret, "data.outputs."+strconv.Itoa(index)+".type").String() == "P2SH" {
		return false, "", errors.New("do not support spending multisig output at the current version")
	}

	scriptPubKey := gjson.Get(ret, "data.outputs."+strconv.Itoa(index)+".script_hex").String()
	return true, scriptPubKey, nil
}
