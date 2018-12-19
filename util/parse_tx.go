package util

import (
	"bytes"
	"encoding/hex"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/blockchain"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/copernet/whc.go/btcjson"
	"github.com/pkg/errors"
)

func DecodeRawTransaction(hexStr string, params *chaincfg.Params) (*btcjson.TxRawDecodeResult, error) {
	// Deserialize the transaction.
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	serializedTx, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, errors.New("hexadecimal data wrong format")
	}
	var mtx wire.MsgTx
	err = mtx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "TX decode failed: " + err.Error(),
		}
	}

	// Create and return the result.
	txReply := btcjson.TxRawDecodeResult{
		Txid:     mtx.TxHash().String(),
		Version:  mtx.Version,
		Locktime: mtx.LockTime,
		Vin:      CreateVinList(&mtx),
		Vout:     CreateVoutList(&mtx, params, nil),
	}
	return &txReply, nil
}

// CreateVinList returns a slice of JSON objects for the inputs of the passed
// transaction.
func CreateVinList(mtx *wire.MsgTx) []btcjson.Vin {
	// Coinbase transactions only have a single txin by definition.
	vinList := make([]btcjson.Vin, len(mtx.TxIn))
	if blockchain.IsCoinBaseTx(mtx) {
		txIn := mtx.TxIn[0]
		vinList[0].Coinbase = hex.EncodeToString(txIn.SignatureScript)
		vinList[0].Sequence = txIn.Sequence
		return vinList
	}

	for i, txIn := range mtx.TxIn {
		// The disassembled string will contain [error] inline
		// if the script doesn't fully parse, so ignore the
		// error here.
		disbuf, _ := txscript.DisasmString(txIn.SignatureScript)

		vinEntry := &vinList[i]
		vinEntry.Txid = txIn.PreviousOutPoint.Hash.String()
		vinEntry.Vout = txIn.PreviousOutPoint.Index
		vinEntry.Sequence = txIn.Sequence
		vinEntry.ScriptSig = &btcjson.ScriptSig{
			Asm: disbuf,
			Hex: hex.EncodeToString(txIn.SignatureScript),
		}
	}

	return vinList
}

// CreateVoutList returns a slice of JSON objects for the outputs of the passed
// transaction.
func CreateVoutList(mtx *wire.MsgTx, chainParams *chaincfg.Params,
	filterAddrMap map[string]struct{}) []btcjson.Vout {
	voutList := make([]btcjson.Vout, 0, len(mtx.TxOut))
	for i, v := range mtx.TxOut {
		// The disassembled string will contain [error] inline if the
		// script doesn't fully parse, so ignore the error here.
		disbuf, _ := txscript.DisasmString(v.PkScript)

		// Ignore the error here since an error means the script
		// couldn't parse and there is no additional information about
		// it anyways.
		scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(
			v.PkScript, chainParams)

		// Encode the addresses while checking if the address passes the
		// filter when needed.
		passesFilter := len(filterAddrMap) == 0
		encodedAddrs := make([]string, len(addrs))
		for j, addr := range addrs {
			encodedAddr := addr.EncodeAddress(true)
			encodedAddrs[j] = encodedAddr

			// No need to check the map again if the filter already
			// passes.
			if passesFilter {
				continue
			}
			if _, exists := filterAddrMap[encodedAddr]; exists {
				passesFilter = true
			}
		}

		if !passesFilter {
			continue
		}

		var vout btcjson.Vout
		vout.N = uint32(i)
		vout.Value = cashutil.Amount(v.Value).ToBCH()
		vout.ScriptPubKey.Addresses = encodedAddrs
		vout.ScriptPubKey.Asm = disbuf
		vout.ScriptPubKey.Hex = hex.EncodeToString(v.PkScript)
		vout.ScriptPubKey.Type = scriptClass.String()
		vout.ScriptPubKey.ReqSigs = int32(reqSigs)

		voutList = append(voutList, vout)
	}

	return voutList
}
