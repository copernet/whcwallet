package api

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bcext/cashutil"
	"github.com/bcext/cashutil/base58"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/whc.go/btcjson"
	"github.com/copernet/whccommon/model"
	"github.com/copernet/whcwallet/config"
	model2 "github.com/copernet/whcwallet/model"
	"github.com/copernet/whcwallet/model/view"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const (
	// 10 satoshi/byte
	feeRate              = 10
	defaultSignatrueSize = 107
	defaultSequence      = 0xffffffff
	dustSatoshi          = 546
	minFeeRate           = 1000 //satoshi/KB
)

func AssembleBCHTx(c *gin.Context) {
	var userTx view.UserTx

	if err := c.ShouldBind(&userTx); err != nil {
		c.JSON(200, apiErrorWithMsg(ErrFormItems, "wrong parameters"))
		return
	}
	if err := validateUserTx(&userTx, c); err != nil {
		c.JSON(200, apiErrorWithMsg(ErrFormItems, "wrong parameters"))
		return
	}

	hex, utxoes, err := doTheBCHAssemble(&userTx, c)
	if err == nil {
		c.JSON(200, apiSuccess(map[string]interface{}{
			"unsigned_tx": *hex,
			"sign_data":   utxoes,
		}))
		return
	}
	c.JSON(200, apiError(ErrAssembleBCHTransaction))
}

func BroadcastBCHTx(c *gin.Context) {
	rawTxHex := c.PostForm("tx")

	resp, err := getNode().BlockchainTransactionBroadcast(rawTxHex)
	if err == nil {
		c.JSON(200, apiSuccess(map[string]interface{}{
			"txHash": resp,
		}))
		return
	}
	c.JSON(200, apiError(ErrSendRawTransaction))
}

func BchHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", strconv.Itoa(defaultPageSize)))
	addr := c.DefaultQuery("addr", "")
	hists, pageCurrent, pageSizeCurrent, total, err := getHistory(addr, page, pageSize)
	if err == nil {
		c.JSON(200, apiSuccess(map[string]interface{}{
			"page":     pageCurrent,
			"pageSize": pageSizeCurrent,
			"total":    total,
			"list":     hists,
		}))
		return
	}
	c.JSON(200, apiError(ErrGetHistoryList))
}

func getScriptPubByAddr(addr string) ([]byte, error) {
	ad, err := cashutil.DecodeAddress(addr, config.GetChainParam())
	if err == nil {
		_, netID, _ := base58.CheckDecode(ad.EncodeAddress(false))
		isP2PKH := chaincfg.IsPubKeyHashAddrID(netID)
		isP2SH := chaincfg.IsScriptHashAddrID(netID)
		var scriptBytes []byte
		if isP2PKH {
			scriptBytes = getP2pkhScript(ad.ScriptAddress())
		}
		if isP2SH {
			scriptBytes = getP2SHScript(ad.ScriptAddress())
		}
		return scriptBytes, nil
	}
	return nil, err
}

func getP2pkhScript(scriptHash []byte) []byte {
	pkScript := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(scriptHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG)

	// ignore error because the specified address is checked
	bs, _ := pkScript.Script()

	return bs
}

func getP2SHScript(scriptHash []byte) []byte {
	p2shScript := txscript.NewScriptBuilder().AddOp(txscript.OP_HASH160).
		AddData(scriptHash).AddOp(txscript.OP_EQUAL)

	// ignore error because the specified address is checked
	bs, _ := p2shScript.Script()

	return bs
}
func validateUserTx(userTx *view.UserTx, c *gin.Context) error {
	if userTx == nil {
		return errors.New("nil userTx")
	}

	if _, err := cashutil.DecodeAddress(userTx.FromAddress, config.GetChainParam()); err != nil {
		return err
	}
	if userTx.ToAddress == nil || len(userTx.ToAddress) == 0 || userTx.ToAmount == nil || len(userTx.ToAmount) == 0 {
		return errors.New("empty userTx's to addresses or amount")
	}
	for _, adr := range userTx.ToAddress {
		if _, err := cashutil.DecodeAddress(adr, config.GetChainParam()); err != nil {
			return err
		}
	}
	if len(userTx.RedeemAddress) > len("undefined") {
		if _, err := cashutil.DecodeAddress(userTx.RedeemAddress, config.GetChainParam()); err != nil {
			return err
		}
	} else {
		userTx.RedeemAddress = ""
	}
	if err := validateAmount(userTx.Fee, false, true); err != nil {
		return err
	}
	for _, amount := range userTx.ToAmount {
		if err := validateAmount(amount, true, false); err != nil {
			return err
		}
	}

	return nil
}
func validateAmount(amount float64, needToCheckDust bool, needToCheckMinFeeRate bool) error {
	if amount <= 0 {
		return errors.New("must great than 0")
	}
	if decimal.NewFromFloat(amount).Exponent() < -8 {
		return errors.New("precision is less than -8")
	}
	if needToCheckDust && decimal.NewFromFloat(amount).Mul(decimal.New(1e8, 0)).IntPart() < dustSatoshi {
		return errors.New("toAmount is less than dust")
	}
	if needToCheckMinFeeRate && decimal.NewFromFloat(amount).Mul(decimal.New(1e8, 0)).IntPart() < minFeeRate {
		return errors.New("fee amount is less than minFeeRate")
	}
	return nil
}
func doTheBCHAssemble(userTx *view.UserTx, c *gin.Context) (*string, []btcjson.PrevTx, error) {
	addr, err := cashutil.DecodeAddress(userTx.FromAddress, config.GetChainParam())
	if err != nil {
		return nil, nil, err
	}
	utxoes, err := GetUtxoElectrumx(addr, decimal.Zero, true)


	exp := decimal.New(1e8, 0)
	//1. assemble outputs' PubkeyScript and Amount, including all the recepients outputs and the change address's output
	var tx wire.MsgTx
	tx.Version = 1
	tx.LockTime = 0
	tx.TxOut = make([]*wire.TxOut, len(userTx.ToAddress)+1)
	for index, to := range userTx.ToAddress {
		pubScriptBytes, _ := getScriptPubByAddr(to)
		tx.TxOut[index] = &wire.TxOut{PkScript: pubScriptBytes, Value: decimal.NewFromFloat(userTx.ToAmount[index]).Mul(exp).IntPart()}
	}
	changeAddress := userTx.FromAddress
	if userTx.RedeemAddress != "" {
		changeAddress = userTx.RedeemAddress
	}
	changePubScriptBytes, _ := getScriptPubByAddr(changeAddress)
	tx.TxOut[len(userTx.ToAddress)] = &wire.TxOut{PkScript: changePubScriptBytes}

	//2. loop the utxoes, collect enough utxoes to fund the totalOutAmount+fee

	totalOutAmount := calTotalOutput(userTx)

	var inputValue int64
	collectUtxos := make([]btcjson.PrevTx, 0)
	successMakeTx := false

	userFeeRate := int(decimal.NewFromFloat(userTx.Fee).Mul(exp).Div(decimal.New(1000, 0)).IntPart())
	for i, utxo := range utxoes {
		utxoValue := decimal.NewFromFloat(utxo.Value).Mul(exp).IntPart()
		if utxoValue <= 0 {
			continue
		}
		inputValue += utxoValue
		collectUtxos = append(collectUtxos, utxo)

		hash, _ := chainhash.NewHashFromStr(utxo.TxID)
		txIn := wire.TxIn{
			PreviousOutPoint: *wire.NewOutPoint(hash, uint32(utxo.Vout)),
			Sequence:         defaultSequence,
		}
		tx.TxIn = append(tx.TxIn, &txIn)

		//option 1: calFee by feeRate:
		fee := (tx.SerializeSize() + defaultSignatrueSize*(i+1)) * userFeeRate
		//option 2: calFee by user defined fixed total fee
		//fee := decimal.NewFromFloat(userTx.Fee).Mul(exp).IntPart()

		offset := inputValue - int64(fee) - totalOutAmount
		if offset >= dustSatoshi {
			tx.TxOut[len(userTx.ToAddress)].Value = int64(offset)
			successMakeTx = true
			break
		}
	}
	if !successMakeTx {
		return nil, nil, errors.New("insufficient funds to make a tx")
	}
	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	txHex := hex.EncodeToString(buf.Bytes())
	return &txHex, collectUtxos, nil

}

func calTotalOutput(userTx *view.UserTx) int64 {
	var ret int64 = 0
	if userTx == nil {
		return ret
	}
	exp := decimal.New(1e8, 0)
	for _, txOutAmount := range userTx.ToAmount {
		ret += decimal.NewFromFloat(txOutAmount).Mul(exp).IntPart()
	}
	return ret
}

type TxHist struct {
	IsReceive         bool            `json:"isReceive"`
	Txid              string          `json:"txid"`
	BlockUTCTimeStamp int64           `json:"blockUTCTimestamp"`
	BalanceDiff       decimal.Decimal `json:"balanceDiff"`
	Confirmations     int64           `json:"confirmations"`
}

func getHistoryHttpContent(addr string, page int, pageSize int) (string, error) {
	url := config.GetBCHAPI() + "/address/" + addr + "/tx?pagesize=" +
		strconv.Itoa(pageSize) + "&page=" + strconv.Itoa(page)

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

func getBCHHistory(addrStr string, page int, pageSize int) ([]model2.HistoryList, int64, error) {
	list, _, _, total, err := getHistory(addrStr, page, pageSize)
	var resuletSet []model2.HistoryList
	for _, v := range list {
		var temp model2.HistoryList
		temp.TxState = string(model.Pending)
		if v.Confirmations > 0 {
			temp.TxState = string(model.Valid)
		}
		temp.TxHash = v.Txid
		temp.TxType = -1      //bch history has no this attribute
		temp.AddressRole = "" //bch history has no this attribute
		temp.BalanceAvailableCreditDebit = v.BalanceDiff
		temp.PropertyName = "BCH"
		temp.PropertyID = 0 //bch history has no this attribute
		temp.BlockTime = v.BlockUTCTimeStamp
		temp.Address = addrStr
		resuletSet = append(resuletSet, temp)
	}

	return resuletSet, total, err
}

func getHistory(addrStr string, page int, pageSize int) ([]TxHist, int64, int64, int64, error) {
	addr, err := cashutil.DecodeAddress(addrStr, config.GetChainParam())
	if err != nil {
		return nil, 1, int64(pageSize), 0, err
	}
	content, err := getHistoryHttpContent(addr.EncodeAddress(false), page, pageSize)
	if err != nil {
		return nil, 1, int64(pageSize), 0, err
	}

	// encounter an error
	remoteErrorNo := gjson.Get(content, "err_no").String()
	if remoteErrorNo != "0" && remoteErrorNo != "1" {
		return nil, 1, int64(pageSize), 0, errors.New(gjson.Get(content, "err_msg").String())
	}

	lists := gjson.Get(content, "data.list").Array()
	ret := make([]TxHist, 0)

	//no history
	if len(lists) == 0 {
		return ret, 1, int64(pageSize), 0, nil
	}

	for _, item := range lists {
		hist := TxHist{
			IsReceive:         item.Get("balance_diff").Int() >= 0,
			Txid:              item.Get("hash").String(),
			BlockUTCTimeStamp: item.Get("block_time").Int(),
			Confirmations:     item.Get("confirmations").Int(),
			BalanceDiff:       decimal.New(item.Get("balance_diff").Int(), -8),
		}
		ret = append(ret, hist)
	}
	p := gjson.Get(content, "data.page").Int()
	total := gjson.Get(content, "data.total_count").Int()
	return ret, p, int64(pageSize), total, nil
}
