package api

import "strconv"

type errCode int

var errNotFound = Response{
	Code:    4004,
	Message: "unknown error",
	Result:  nil,
}

// consensus error
const (
	ErrUnSupportTxType errCode = 2000 + iota
	ErrUnSupportTxVersion
	ErrUnSupportEcosystem
	ErrUnSupportCrowdSaleSP
	ErrUnSupportPrecision
)

// RPC error
const (
	ErrCreatePayload errCode = 2100 + iota
	ErrCreateRawTxInput
	ErrCreateRawTxReference
	ErrCreateRawTxOpReturn
	ErrCreateRawTxChange
	ErrSendRawTransaction
	ErrWhcDecodeTransaction
	ErrWhcGetTransaction
	ErrAssembleBCHTransaction
)

// User input error
const (
	ErrFormItems errCode = 2200 + iota
	ErrConvertFloat64
	ErrEmptyAddressList
	ErrGetTxType
	ErrIncorrectAddress
	ErrConvertInt
	ErrIncorrectAmount
	ErrEmptyQueryParam
	ErrCategoryNotFound
	ErrHexStringFormat
	ErrHash256Format
	ErrIncorrectPropertyID
	ErrIncorrectPropertyName
	ErrIncorrectPropertyQueryByKeyword
	ErrExceedMaxAddressRequestLimit
	ErrExceedMaxNulldataLimit
	ErrEmptyDeviceName
	ErrParseVersionString
	ErrHightestVersionString
)

// database error
const (
	ErrGetHistoryList errCode = 2300 + iota
	ErrGetHistoryListCount
	ErrGetHistoryDetail
	ErrEmptyHistoryDetail
	ErrEmptyHistorySpDetail
	ErrGetPropertyByAddressIssuer
	ErrGetPropertyByAddressIssuerCount
	ErrGetPropertyByAddress
	ErrGetPropertyByID
	ErrGetPropertyByName
	ErrListProperties
	ErrListPropertiesCount
	ErrGetAllCrowdSale
	ErrGetActiveCrowdSale
	ErrGetActiveCrowdSaleCount
	ErrQueryTotal
	ErrQueryTransactions
	ErrGetBalanceFromRedis
	ErrGetBalanceFromDatabase
	ErrEmptyBalance
	ErrGetPurchasedCrowdSaleList
	ErrGetPurchasedCrowdSaleCount
	ErrGetNewestVersionFailed
	ErrGetNotification
)

// account error
const (
	ErrChallenge errCode = 2400 + iota
	ErrCreateWallet
	ErrUpdateWallet
	ErrLogin
	ErrNotUUID
	ErrVerify
	ErrMfaToken
	ErrInsufficientBalance
)

// server or handle error
const (
	ErrInsertBCH errCode = 2500 + iota
	ErrInsertWhc
	ErrCanNotGetUtxo
	ErrCanNotGetInputs
	ErrDecodeRawTransaction
	ErrTxDeserialize
	ErrGetBCHBalance
	ErrVersionSetting
	ErrAmountString
)

const (
	ErrTransactionGet errCode = 2600 + iota
)

var errMapping = map[errCode]*Response{
	// consensus error
	ErrUnSupportTxType:      {Code: ErrUnSupportTxType, Message: "Unsupported tx type"},
	ErrUnSupportTxVersion:   {Code: ErrUnSupportTxVersion, Message: "Unsupported tx version"},
	ErrUnSupportEcosystem:   {Code: ErrUnSupportEcosystem, Message: "Unsupported ecosystem"},
	ErrUnSupportCrowdSaleSP: {Code: ErrUnSupportCrowdSaleSP, Message: "Unsupported crowdsale smart property"},
	ErrUnSupportPrecision:   {Code: ErrUnSupportPrecision, Message: "Unsupported property precision"},

	// RPC error
	ErrCreatePayload:          {Code: ErrCreatePayload, Message: "Create wormhole transaction payload error"},
	ErrCreateRawTxInput:       {Code: ErrCreateRawTxInput, Message: "RPC call for whc_createrawtx_input error"},
	ErrCreateRawTxReference:   {Code: ErrCreateRawTxReference, Message: "RPC call for whc_createrawtx_reference error"},
	ErrCreateRawTxOpReturn:    {Code: ErrCreateRawTxOpReturn, Message: "RPC call for whc_createrawtx_opreturn error"},
	ErrCreateRawTxChange:      {Code: ErrCreateRawTxChange, Message: "RPC call for whc_createrawtx_change error"},
	ErrSendRawTransaction:     {Code: ErrSendRawTransaction, Message: "RPC call for sendrawtransaction error"},
	ErrWhcDecodeTransaction:   {Code: ErrWhcDecodeTransaction, Message: "RPC call for whc_decodetransaction error"},
	ErrWhcGetTransaction:      {Code: ErrWhcDecodeTransaction, Message: "RPC call for whc_gettransaction fro wormhole transaction error"},
	ErrAssembleBCHTransaction: {Code: ErrAssembleBCHTransaction, Message: "Assemble BCH transaction error, available UTXOs is not enough currently or System error"},

	// User input error
	ErrFormItems:                       {Code: ErrFormItems, Message: "Please check form input"},
	ErrConvertFloat64:                  {Code: ErrConvertFloat64, Message: "Float number incorrect"},
	ErrEmptyAddressList:                {Code: ErrEmptyAddressList, Message: "Address list empty"},
	ErrGetTxType:                       {Code: ErrGetTxType, Message: "Get txtype failed"},
	ErrIncorrectAddress:                {Code: ErrIncorrectAddress, Message: "Unrecognized address format"},
	ErrConvertInt:                      {Code: ErrConvertInt, Message: "Integer number incorrect"},
	ErrIncorrectAmount:                 {Code: ErrIncorrectAmount, Message: "Bitcoin cash or wormhole number incorrect"},
	ErrEmptyQueryParam:                 {Code: ErrEmptyQueryParam, Message: "Empty query parameter"},
	ErrCategoryNotFound:                {Code: ErrCategoryNotFound, Message: "Category not found"},
	ErrHexStringFormat:                 {Code: ErrHexStringFormat, Message: "Hexadecimal string misform"},
	ErrHash256Format:                   {Code: ErrHash256Format, Message: "Sha256 hash string misform"},
	ErrIncorrectPropertyID:             {Code: ErrIncorrectPropertyID, Message: "Incorrect property id"},
	ErrIncorrectPropertyName:           {Code: ErrIncorrectPropertyName, Message: "Incorrect property name"},
	ErrIncorrectPropertyQueryByKeyword: {Code: ErrIncorrectPropertyQueryByKeyword, Message: "Incorrect property query keyword"},
	ErrExceedMaxAddressRequestLimit:    {Code: ErrExceedMaxAddressRequestLimit, Message: "So many addresses requested(max:" + strconv.Itoa(maxRequestAddressList) + ")"},
	ErrEmptyDeviceName:                 {Code: ErrEmptyDeviceName, Message: "Please input device fields"},
	ErrParseVersionString:              {Code: ErrParseVersionString, Message: "Make sure your version string valid"},
	ErrHightestVersionString:           {Code: ErrHightestVersionString, Message: "Invalid version higher than the server given"},

	// database
	ErrGetHistoryList:                  {Code: ErrGetHistoryList, Message: "Get history list failed"},
	ErrGetHistoryListCount:             {Code: ErrGetHistoryListCount, Message: "Get history list total count failed"},
	ErrGetHistoryDetail:                {Code: ErrGetHistoryDetail, Message: "Get history detail failed"},
	ErrEmptyHistoryDetail:              {Code: ErrEmptyHistoryDetail, Message: "Get history detail empty"},
	ErrEmptyHistorySpDetail:            {Code: ErrEmptyHistorySpDetail, Message: "Get history smart property detail empty"},
	ErrGetPropertyByAddressIssuer:      {Code: ErrGetPropertyByAddressIssuer, Message: "Get property list issued by the address failed"},
	ErrGetPropertyByAddressIssuerCount: {Code: ErrGetPropertyByAddressIssuerCount, Message: "Get property list issued number by the address failed"},
	ErrGetPropertyByAddress:            {Code: ErrGetPropertyByAddress, Message: "Get property list by address failed"},
	ErrGetPropertyByID:                 {Code: ErrGetPropertyByID, Message: "Get property data by ID failed"},
	ErrGetPropertyByName:               {Code: ErrGetPropertyByName, Message: "Get property data by Name failed"},
	ErrListProperties:                  {Code: ErrListProperties, Message: "Get properties list failed"},
	ErrListPropertiesCount:             {Code: ErrListPropertiesCount, Message: "Get property list total number failed"},
	ErrGetAllCrowdSale:                 {Code: ErrGetAllCrowdSale, Message: "Get all crowdsale failed"},
	ErrGetActiveCrowdSale:              {Code: ErrGetActiveCrowdSale, Message: "Get actives crowdsale failed"},
	ErrGetActiveCrowdSaleCount:         {Code: ErrGetActiveCrowdSaleCount, Message: "Get actives crowdsale list total count failed"},
	ErrQueryTotal:                      {Code: ErrQueryTotal, Message: "Query total failed"},
	ErrQueryTransactions:               {Code: ErrQueryTransactions, Message: "Query transaction failed"},
	ErrGetBalanceFromRedis:             {Code: ErrGetBalanceFromRedis, Message: "Get balance from cache failed"},
	ErrGetBalanceFromDatabase:          {Code: ErrGetBalanceFromDatabase, Message: "Get balance failed"},
	ErrEmptyBalance:                    {Code: ErrEmptyBalance, Message: "Balance of the specified address empty"},
	ErrGetPurchasedCrowdSaleList:       {Code: ErrGetPurchasedCrowdSaleList, Message: "Get purchased crowdsale list failed"},
	ErrGetPurchasedCrowdSaleCount:      {Code: ErrGetPurchasedCrowdSaleCount, Message: "Get purchased crowdsale count failed"},
	ErrGetNewestVersionFailed:          {Code: ErrGetNewestVersionFailed, Message: "Get newest version failed"},
	ErrGetNotification:                 {Code: ErrGetNotification, Message: "Get notification failed"},

	// account error
	ErrChallenge:           {Code: ErrChallenge, Message: "Account challenge error"},
	ErrCreateWallet:        {Code: ErrCreateWallet, Message: "Account create wallet error"},
	ErrUpdateWallet:        {Code: ErrUpdateWallet, Message: "Account update wallet failed"},
	ErrLogin:               {Code: ErrLogin, Message: "Account login failed"},
	ErrNotUUID:             {Code: ErrNotUUID, Message: "Account not found the specified uuid"},
	ErrVerify:              {Code: ErrVerify, Message: "Account verify error"},
	ErrMfaToken:            {Code: ErrMfaToken, Message: "Account has set mfa_token"},
	ErrInsufficientBalance: {Code: ErrInsufficientBalance, Message: "Account has insufficient balance fro creating transaction"},

	// server or handle error
	ErrInsertBCH:            {Code: ErrInsertBCH, Message: "Insert bitcoin-cash transaction failed"},
	ErrInsertWhc:            {Code: ErrInsertWhc, Message: "Insert wormhole transaction failed"},
	ErrCanNotGetUtxo:        {Code: ErrCanNotGetUtxo, Message: "Can not fetch utxo"},
	ErrCanNotGetInputs:      {Code: ErrCanNotGetInputs, Message: "Can not fetch transaction's inputs information"},
	ErrDecodeRawTransaction: {Code: ErrDecodeRawTransaction, Message: "Decode raw transaction failed"},
	ErrTxDeserialize:        {Code: ErrTxDeserialize, Message: "Deserialize transaction failed"},
	ErrGetBCHBalance:        {Code: ErrGetBCHBalance, Message: "Fetch bitcoin cash balance from electrum server failed"},
	ErrVersionSetting:       {Code: ErrVersionSetting, Message: "Server data malformed"},
	ErrAmountString:         {Code: ErrAmountString, Message: "Server amount data malformed"},

	// electrum error
	ErrTransactionGet: {Code: ErrTransactionGet, Message: "Get transaction detail failed"},
}

func apiErrorWithMsg(code errCode, msg string) *Response {
	r, ok := errMapping[code]
	if !ok {
		return &errNotFound
	}

	if msg != "" {
		// copy the origin response
		customMsg := *r
		customMsg.Message = msg
		return &customMsg
	}

	return r
}

func apiError(code errCode) *Response {
	r, ok := errMapping[code]
	if !ok {
		return &errNotFound
	}

	return r
}

func apiSuccess(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "",
		Result:  data,
	}
}
