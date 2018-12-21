##### `/property/listowners/:id`

Ruturn the owner information and amount of property，address，balance_available and status will be returned.

**Argument**

- pageNo `string` required:page number (can be null)
- pageSize `string` required:page size(can be null)

##### Result

`list`  owner list

- `Address`   account address

- `balance_available`  owner balance
- `Status` owner status

`total`  eligible results amount

##### Example

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"pageSize":"5","pageNo":"1"} http://127.0.0.1:9999/property/listowners/459
//Result
{"result": {"list": [{
"Address": "bchtest:qrhncn6hgkhljqg4fuq4px5qg74sjz9fqqkg3h8g6e",
"balance_available": "0","Status": 0}],"total": 1}}

```

##### `/getunsigned/185`

Return the unsigned transaction of freeze token ,sign_data and unsigned_tx will be returned.

| Protocol | Method | API              |
| -------- | ------ | ---------------- |
| HTTP     | POST   | /getunsigned/185 |

**Argument**

- transaction_version `string`  required:transaction version number
- pubkey `string`required：public key
- fee `string` required：transaction fee
- transaction_from `string` required：transaction issuer address
- property_id `string` required: property number
- frozen_address `string` required:addresses that need to be frozen

**Result**

`sign_data`

- `txid`  transaction hash
- `vout`  past output index
- `scriptPubKey`  lock script
- `value`   value

`unsigned_tx`：serialization of unsigned transactions

**Example**

```go
//request
curl -i -H 'Content-Type: application/json' -X POST -d{"transaction_version":0,"fee":"0.001,"transaction_from":"bchtest:qrwejrccj3gfhtf025yt27fjhup09lxkfc0h0yl7ev","property_id":"459","amount":"1","frozen_address":"bchtest:qqu9lh4jpc05p59pfhu9amyv9uvder8j3sa2up95v","pubkey":""} http://127.0.0.1:9999/getunsigned/185
//result
{"result":{"sign_data":[{"txid":"131b64d4f5903c60a48d2a336f03b80b6ef786b3b538ba0e3ccfd546a3444606","vout":0,"scriptPubKey":"76a914dd990f1894509bad2f5508b57932bf02f2fcd64e88ac","value":0.8097907}],"unsigned_tx":"0200000001064644a346d5cf3c0eba38b5b386f76e0bb8036f332a8da4603c90f5d4641b130000000000ffffffff027639d304000000001976a914dd990f1894509bad2f5508b57932bf02f2fcd64e88ac0000000000000000496a4708776863000000b9000001cb0000000005f5e100626368746573743a717175396c68346a706330357035397066687539616d7976397576646572386a337361327570393576730000000000"}}
```

##### `/getunsigned/186`

Return the unsigned transaction of unfreeze token ,sign_data and unsigned_tx will be returned.

| Protocol | Method | API              |
| -------- | ------ | ---------------- |
| HTTP     | POST   | /getunsigned/186 |

**Argument **

- transaction_version `string` required:transaction version number
- pubkey `string` :public key
- fee `string` required:transaction fee
- transaction_from `string` required: transaction issuer address
- property_id `string` required:property type
- amount `string` required: unfreeze token  amount
- frozen_address  `string` required:freeze address

**Result**

`sign_data`

- `txid`  transaction hash
- `vout` past output index
- `scriptPubKey`  lock script
- `value`    value

`unsigned_tx`：serialization of unsigned transactions

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"transaction_version":0,"fee":"0.001","transaction_from":"bchtest:qrwejrccj3gfhtf025yt27fjhup09lxkfc0h0yl7ev","property_id":"459","amount":"1","frozen_address":"bchtest:qqu9lh4jpc05p59pfhu9amyv9uvder8j3sa2up95v","pubkey":""} http://127.0.0.1:9999/getunsigned/186
//Result
{"result":{"sign_data":[{"txid":"131b64d4f5903c60a48d2a336f03b80b6ef786b3b538ba0e3ccfd546a3444606","vout":0,"scriptPubKey":"76a914dd990f1894509bad2f5508b57932bf02f2fcd64e88ac","value":0.8097907}],"unsigned_tx":"0200000001064644a346d5cf3c0eba38b5b386f76e0bb8036f332a8da4603c90f5d4641b130000000000ffffffff027639d304000000001976a914dd990f1894509bad2f5508b57932bf02f2fcd64e88ac0000000000000000496a4708776863000000ba000001cb0000000005f5e100626368746573743a717175396c68346a706330357035397066687539616d7976397576646572386a337361327570393576730000000000"}}
```

##### `/device/update`

Return the newest version information for device

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | GET    | /device/update |

**Argument**

None

**Result**

- `version`:newest version number
- `download`:download link
- `description`：newest version description

**Example**

```go 
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/device/update/android
//Result
{"result":{"version":"0.1.0beta","versionCode":4,"download":"https://wallet.wormhole.cash","description":"1、support look up you seed.\\n 2、fix some bugs"}}
```

##### `/notify`

Return the transaction information of the address within a certain period of block time

| Protocol | Method | API     |
| :------- | ------ | ------- |
| HTTP     | GET    | /notify |

**Argument**

- address `string` required: address to which notifications need to be queried
- from `string` required:begin block time
- to `string` required: end block time

**Result**

- `address`:
- `property_id`：property id
- `property_name`：property name
- `protocol`：protocol
- `address_role`：address role
- **`balance_available_credit_debit`：**token balance 
- `balance_frozen_credit_debit`：froze token count
- `block_time`：block time
- `tx_type_name`：transaction type name

**Example**

```go 
//Request
curl -i -H 'Content-Type: application/json' -X POST -d '{"address": "bchtest:qq6qag6m"2f3djnv7ddk3lpk", "from":"1","to":"1541141668"}'  http://127.0.0.1:9999/notify/
//Result
{"result":{"list":[{"address":"bchtest:qp437dgcppz2u5uc93k329tgvv804djaxylvp0v5e2","property_id":445,"property_name":"GBDCT","protocol":"wormhole","address_role":"recipient","balance_available_credit_debit":"500","balance_reserved_credit_debit":null,"balance_accepted_credit_debit":null,"balance_frozen_credit_debit":"0","block_time":1540956811,"tx_type_name":"Simple Send"},{"address":"bchtest:qzwgwvyal5fmqyz3z2av0vq5wrelxlp5xv02r8d6ku","property_id":445,"property_name":"GBDCT","protocol":"wormhole","address_role":"sender","balance_available_credit_debit":"-500","balance_reserved_credit_debit":null,"balance_accepted_credit_debit":null,"balance_frozen_credit_debit":"0","block_time":1540956811,"tx_type_name":"Simple Send"}],"total":2}}
```

##### `/bch/hists`

Return transaction records of BCH.

| Protocol | Method | API        |
| -------- | ------ | ---------- |
| HTTP     | GET    | /bch/hists |

**Argument**

- addr`string`required:wallet address
- page`number`: page number to search(default 1)

**Result**

- `total`：transaction records number
- `pageSize`：page size
- `page`：current page

- `list`：transaction record list
  - `  isReceive`:income or expenditure
  - `txid`:transaction hash
  - `blockUTCTimestamp`:block time is in accordance with UTC time zone
  - `balanceDiff` : account balance changes as transaction
  - `confirmations`: transaction confirmation

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"addr":"bchtest:qryk7m37zfe0nycl4227hw2deardw9tlvswkj5aw8a","page":"1"}' http://127.0.0.1:9999/bch/hists

//Result
{"result": {"list": [ { "isReceive": false,"txid":"70b6803642f07d468986fe80d3fa3d7492421e2e5a6acca1f8813f46044c21df","blockUTCTimestamp": 1539073980, "balanceDiff": "-1.0647728", "confirmations": 230}, {"isReceive": true,"txid": "e406f5bae8aafe588f329971495139b15740604138790711f8c03687e7cf3f5d","blockUTCTimestamp": 1539072656, "balanceDiff": "1.0647728", "confirmations": 232 } ],"total": 2,"page": 1, "pageSize": 50}}
```

##### `transaction/fee`

Return transaction costs.

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | transaction/fee |

**Argument**

None

**Result**

- `Fast`:fast rate
- `Normal`:normal rate
- `Slow`:slow rate

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X p://127.0.0.1:9999/transaction/fee
//Result
{"result": {"Fast": 0.00026221,"Normal": 0.0002017,"Slow": 0.00014119}}
```

##### `/user/wallet/challenge`

Return salt data.

| Protocol | Method | API                    |
| -------- | ------ | ---------------------- |
| HTTP     | GET    | /user/wallet/challenge |

**Argument**

- uuid`string` required:unique identification

**Result**

- `pow_challenge`pow challenge
- `salt`: salt data
- `challenge` challenge

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/user/wallet/challenge
//Result
{"result":{"salt":"53e92f6214df9d7fef8a75557f685e637eb875e6b626177b65d689691b39a00a","pow_challenge":"96157b0a626c53fdd970dd5536b4760d","challenge":"0c8eb4307abed4c0a142388c71eb6853"}}
```

##### `/user/wallet/create`

Create wallet and update session public key

| Protocol | Method | API                 |
| -------- | ------ | ------------------- |
| HTTP     | POST   | /user/wallet/create |

**Argument**

- uuid`string`required: unique user ID
- Email`string`required:email
- nonce`string`:nonce
- public_key`string`:public key
- wallet`string`:encryption of wallets with a given key
- recaptcha_response_field`string`: recaptcha response field

`/user/wallet/update`

Update wallet info

| Protocol | Method | Api                 |
| -------- | ------ | ------------------- |
| HTTP     | POST   | /user/wallet/update |

**Argument**

- Uuid`string`required: unique user ID
- wallet`string` required:encryption of wallets with a given key
- Signature`string`required:signature
- email`string` required:email

##### `/user/wallet/login`

For login account

| Protocol | Method | API                |
| -------- | ------ | ------------------ |
| HTTP     | POST   | /user/wallet/login |

**Argument**

- nonce`string` reqiured:nonce
- public_key`string` required:public key
- uuid`string` required: unique user  ID
- Email`string`:email

##### `/user/wallet/verify`

The api is used for web ajax email unique or phone unique.

| Protocol | Method | API                 |
| -------- | ------ | ------------------- |
| HTTP     | POST   | /user/wallet/verify |

**Argument**

- email`string` required:email

##### `/user/wallet/newMfa`

The api is used for create new mfa

| Protocol | Method | API                 |
| -------- | ------ | ------------------- |
| HTTP     | POST   | /user/wallet/newMfa |

**Argument**

- uuid`string`required:unique user  ID

##### `/ws/balance`

The api is used for  address operation.

| Protocol | Method | API         |
| -------- | ------ | ----------- |
| HTTP     | GET    | /ws/balance |

**Argument**

- None

**Example**

```GO 
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/ws/balance?type=add&addresses="bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa"
```

##### `/bch/assemble`

Return the payment transaction information.

| Protocol | Method | API           |
| -------- | ------ | ------------- |
| HTTP     | POST   | /bch/assemble |

**Argument**

- transaction_from`string`required:transaction from address
- address[]`array`required:receive address
- amount[]`array` required:receive amount
- fee`string` required:transaction fee
- redeem_address`string`required:change address

**Result**

`sign_data`

- `txid` : transaction hash
- `vout`  :past output index
- `scriptPubKey`  :lock script
- `value`  :value

`unsigned_tx`：serialization of unsigned transactions

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"transaction_from":"bchtest:qpvnw78dtuglcl4l6sv9rgysz39eywd9vgdfzggac0","address[]":"bchtest:qpvnw78dtuglcl4l6sv9rgysz39eywd9vgdfzggac0","amount[]":"0.01","fee":"0.0001","redeem_address":"bchtest:qpvnw78dtuglcl4l6sv9rgysz39eywd9vgdfzggac0"}' http://127.0.0.1:9999/bch/assemble
//Result
{"result":{"sign_data":[{"txid":"647a12ad9f7ef4aa53edfe18af062978f4712c56d199d53144e8cf031a3f2d54","vout":0,"scriptPubKey":"76a914593778ed5f11fc7ebfd41851a090144b9239a56288ac","value":0.99999328}],"unsigned_tx":"0100000001542d3f1a03cfe84431d599d1562c71f4782906af18feed53aaf47e9fad127a640000000000ffffffff0440420f000000000017a914c50a16cff8336e27cae53013cb459b85786c74fa8780841e00000000001976a9147a1ef4db3f2211e2362981b33c750558a8bf2e6d88ac40420f000000000017a914c50a16cff8336e27cae53013cb459b85786c74fa870ccab805000000001976a914593778ed5f11fc7ebfd41851a090144b9239a56288ac00000000"}}
```

##### `/bch/broadcast`

To broadcast the signed bch transactions.

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | POST   | /bch/broadcast |

**Argument**

- tx `string` required:the signed transactions hash value

**Result**

- `txHash`:Return the hash value after broadcasting

**Example**

```go 
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"tx":"0100000001f4818c6a24c3dd5482979660cd97539fab665d28888ff709a6e3d0e841edc2320000000000ffffffff0340420f000000000017a914c50a16cff8336e27cae53013cb459b85786c74fa8780841e00000000001976a9147a1ef4db3f2211e2362981b33c750558a8bf2e6d88ac6c1cef0c000000001976a914593778ed5f11fc7ebfd41851a090144b9239a56288ac00000000"}' http://127.0.0.1:9999/bch/broadcast

//Result
{ "result": { "txHash": "ebfe99a35ab72fb0dbee7bcecd61ac26e636a52984bd2224748856145b737d84"
 }}
```

##### `/balance/bch/addresses`

Returns the address balance,the confirmed and unconfirmed will be returned.

| Portocol | Method | API                    |
| -------- | ------ | ---------------------- |
| HTTP     | POST   | /balance/bch/addresses |

**Argument**

- Address `string` required:destination address

**Result**

- `address`: current address
- `confirmed`:confirmed balance
- `unconfirmed`:unconfirmed balance

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg"}' http://127.0.0.1:9999/balance/bch/addresses
//Result
{"result":[{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","confirmed":0.31466998,"unconfirmed":0}]}
```

##### `/crowdsale/list/active`

Returns all active crowdsourcing lists in the current system.

| Protocol | Method | API                    |
| -------- | ------ | ---------------------- |
| HTTP     | GET    | /crowdsale/list/active |

**Argument**

- pageSize`string` required:page size
- pageNo`string`required:page number
- keyword `string`:property number or name(can be null)

**Result**

- `list`
  - `active`crowdsale status
  - `addedissuertokens`added issuer tokens
  - `amountraised` whc amount
  - `category`category
  - `closetx`  close transaction hash
  - `created`creation time
  - `creationtxid` transaction hash  creation
  - `data`reamrk
  - `deadline` crowdsale deadline
  - `earlybonus`weekly early bird bonus percentage
  - `fixedissuance` property type
  - `freezingenabled` freezing enabled
  - `issuer`issuer address
  - `managedissuance`property type
  - `name`crowdsale name
  - `participanttransactions`participant transactions
  - `precision`precision
  - `propertyid`crowdsale id
  - `propertyiddesired`participate in the crowdsale token id
  - `purchasedtokens`token count
  - `starttime`start time
  - `subcategory` subcategory
  - `tokensissued` issued tokens
  - `tokensperunit` whc and tokens proportion
  - `totaltokens` total tokens
  - `url`url
- `total`  eligible results amount

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/crowdsale/list/active?pageSize=5&pageNo=1&keyword=22
//Result
{"result":{"list":[{"active":true,"addedissuertokens":"0.0000000","amountraised":"64.77041915","category":"","closetx":null,"created":1532970004,"creationtxid":"5cba679e754f038062da9becf7f2f41d5cd093cde8bd306179611d27a02f08cd","data":"hi","deadline":1596127632,"earlybonus":10,"fixedissuance":false,"freezingenabled":null,"issuer":"bchtest:qz3fgledq0tgl0ry6pn0c5nmufspfrr8aqsyuc39yl","managedissuance":false,"name":"zc_crow2","participanttransactions":null,"precision":"7","propertyid":27,"propertyiddesired":1,"purchasedtokens":665.9166694,"starttime":1532970004,"subcategory":"","tokensissued":"100000","tokensperunit":"1","totaltokens":"100000","url":"www.zc"}],"total":1}}
```

##### `crowdsale/purchase/times/id/:id`

Return number of times that crowdsourcing was purchased.

| Protocol | Method | API                              |
| -------- | ------ | -------------------------------- |
| HTTP     | GET    | crowdsale/purchase/times/id/:id` |

**Argument**

None

**Result**

- `total`:purchases quantity

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/crowdsale/purchase/times/id/7
//Result
{"result":{"total":3}}
```

##### `/crowdsale/purchase/record/id/:id`

Return to the crowdsourcing participation record

| Procotol | Method | API                               |
| -------- | ------ | --------------------------------- |
| HTTP     | GET    | /crowdsale/purchase/record/id/:id |

**Argument**

- pageNo `string` required:page number(can be null)
- pageSize `string` required:page size(can be null)

**Result**

- `list`
  - `fee`transaction fee
  - `txid`transaction hash
  - `type`transaction type
  - `block`block in
  - `valid`effectiveness
  - `amount`total cost
  - `version`version
  - `type_int`type number
  - `blockhash`block hash
  - `blocktime`block time
  - `precision`precision
  - `propertyid`property id
  - `issuertokens`  issuer tokens
  - `confirmations` confirmation count
  - `actualinvested`actual invested
  - `sendingaddress`transaction issuer address
  - `purchasedtokens`total purchase tokens
  - `referenceaddress`transaction receive address
  - `purchasedpropertyid`purchase property id
  - `purchasedpropertyname`purchase property name
  - `purchasedpropertyprecision`purchase property  precision
- `total`eligible results amount

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/crowdsale/purchase/record/id/7?pageSize=1&pageNo=1
//Result
{"result":{"list":[{"fee":"5665","txid":"09e9c431f3ce3a62b2625da0f4de7d55abed11a9d97aaf6e5af68c64b44eed22","type":"Crowdsale Purchase","block":1249368,"valid":true,"amount":"1.01000001","ismine":false,"version":0,"subsends":null,"type_int":1,"blockhash":"00000000000001a82df433dc8cd94260f64db49f9bdfe8410efce3f9d4b84665","blocktime":1532969235,"precision":"8","propertyid":7,"issuertokens":"0","confirmations":22154,"actualinvested":"1.00995317","sendingaddress":"bchtest:qqwmlaph6l9k4xvs6j3u4mg0ynhvdm3tr5qzn0gmkg","purchasedtokens":"0.11612","referenceaddress":"bchtest:qz04wg2jj75x34tge2v8w0l6r0repfcvcygv3t7sg5","purchasedpropertyid":7,"purchasedpropertyname":"crowsale_5","purchasedpropertyprecision":"5"}],"total":15}}
```

##### `/getunsigned/0`

Return the unsigned transaction of simple send ,sign_data and unsigned_tx will be returned.

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | POST   | /getunsigned/0 |

**Argument**

- transaction_version`string`required:version number of transaction
- pubkey`string`:public key
- fee`string` required:the cost of completing the transaction
- transaction_from`string` required:the issuer address of the transaction
- transaction_to`string` required:the receive address of the transaction
- currency_identifier`string` required:currency identifier
- amount_to_transfer`string` required:amount of simple send transaction
- redeem_address`string`:change address
- reference_amount`string` :incidental bch count

**Result**

- `unsigned_tx`：serialization of unsigned transactions

- ``sign_data`
  - `txid`  transaction hash
  - `vout`  index of past output
  - `scriptPubKey`  lock script
  - `value`    value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","transaction_to":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","currency_identifier":"1","amount_to_transfer":"0.003","redeem_address":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","reference_amount":"0.00001"}' http://127.0.0.1:9999/getunsigned/0

//Result
{"result":{"sign_data":[{"txid":"a0b2b967115401f4a558f4c44ff88340193a6e2e124c905ff1c2023bd8b6160a","vout":0,"scriptPubKey":"76a9141522a025f2365eebee65cd8a8b8a38180dbcd59588ac","value":0.1}],"unsigned_tx":"02000000010a16b6d83b02c2f15f904c122e6e3a194083f84fc4f458a5f401541167b9b2a00000000000ffffffff038e889800000000001976a914340ea35b62922e03d10767bd5b4f70424545b29b88ac0000000000000000166a1408776863000000000000000100000000000493e0e8030000000000001976a914340ea35b62922e03d10767bd5b4f70424545b29b88ac00000000"}}
```

##### `/getunsigned/1`

Return the unsigned transaction of participate crowdsale ,sign_data and unsigned_tx will be returned.

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | POST   | /getunsigned/1 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:participants address
- transaction_to`string` required:crowdsourcing issuer address
- amount_to_transfer`string` required:transaction amount of participate crowdsale 
- redeem_address`string`:change address
- reference_amount`string` :incidental bch count

**Result**

- `unsigned_tx`：serialization of unsigned transactions

- ``sign_data`
  - `txid`  transaction hash
  - `vout`   past output index
  - `scriptPubKey`  lock script
  - `value`    value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","transaction_to":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","currency_identifier":"1","amount_to_transfer":"0.003","redeem_address":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","reference_amount":"0.00001"}' http://127.0.0.1:9999/getunsigned/1

//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff0381cea435000000001976a914340ea35b62922e03d10767bd5b4f70424545b29b88ac0000000000000000166a1408776863000000010000000100000000000186a0e8030000000000001976a914678fe8f46d32484572ab0961ddd5a3fc9c73b02c88ac00000000"}}
```

##### `/getunsigned/3`

Return the unsigned transaction of send to owners,sign_data and unsigned_tx will be returned.

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | POST   | /getunsigned/3 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required: transaction fee
- transaction_from`string` required:issuer address
- redeem_address`string` :change address
- Propertyid`string`required:property id
- Amount`string`required:transaction amount
- Distributionproperty`string`required:target property id

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout`  past output index
  - `scriptPubKey`  lock script
  - `value`   value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","propertyid":"1","amount":"0.987","redeem_address":"bchtest:qrngs64rn42cana8nxqakyv0yw8ceg3sss0dfp4cp0","distributionproperty":"23"}' http://127.0.0.1:9999/getunsigned/3
//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff03efcea435000000001976a914e6886aa39d558ecfa79981db118f238f8ca2308488ac00000000000000000b6a09087768630000000401e8030000000000001976a914678fe8f46d32484572ab0961ddd5a3fc9c73b02c88ac00000000"}}
```

##### `/getunsigned/4`

Return the unsigned transaction of send all handler ,sign_data and unsigned_tx will be returned.

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | POST   | /getunsigned/4 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:transaction issuer address 
- transaction_to`string` required:transaction receive address  
- redeem_address`string` :change address
- reference_amount`string` :incidental bch count
- ecosystem`string`required:ocation system

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout` past output index
  - `scriptPubKey`  lock script
  - `value`   value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","transaction_to":"bchtest:qpncl685d5eys3tj4vykrhw4507fcuas9scm4qd70e","redeem_address":"bchtest:qrngs64rn42cana8nxqakyv0yw8ceg3sss0dfp4cp0","reference_amount":"0.00001","ecosystem":"1"}' http://127.0.0.1:9999/getunsigned/4

//Result
{result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff03efcea435000000001976a914e6886aa39d558ecfa79981db118f238f8ca2308488ac00000000000000000b6a09087768630000000401e8030000000000001976a914678fe8f46d32484572ab0961ddd5a3fc9c73b02c88ac00000000"}}
```

##### `/getunsigned/50`

Return the  unsigned transaction of creating fixed issuance property ,sign_data and unsigned_tx will be returned.   

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/50 |

**Argument**

- transaction_version`string`required: transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required: transaction  issuer address 
- ecosystem`string`required:location system
- precision`string` required:precision
- previous_property_id `string`required:previous property id
- property_category `string` required:property category
- property_subcategory`string`required:property subcategory
- property_name`string`required:property name
- property_url`string` required:property url
- property_data `string` required:property data
- number_properties`string`required: total issue amount

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout` past output index
  - `scriptPubKey`  
  - `value`  

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","ecosystem":"1","precision":"5","previous_property_id":"","property_category":"hello golang","property_subcategory":"whcwallet","property_name":"golang","property_url":"https://cash.wormhole.org","property_data":"hello world","number_properties":"12238432.9882"}' http://127.0.0.1:9999/getunsigned/4

//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff02edd0a435000000001976a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac00000000000000005e6a4c5b08776863000000320100050000000068656c6c6f20676f6c616e670077686377616c6c657400676f6c616e670068747470733a2f2f636173682e776f726d686f6c652e6f72670068656c6c6f20776f726c64000000011cf2bebe0400000000"}}
```

##### `/getunsigned/51`

Return the unsigned transaction of creating crowdSale issuance ,sign_data and unsigned_tx will be returned.   

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/51 |

**Argument**

- transaction_version`string`required:transaction version number 
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:  transaction issuer address 
- ecosystem`string`required:location system
- precision`string` required:precision
- previous_property_id `string`required:previous property id
- property_category `string` required:property category
- property_subcategory`string`required:property subcategory
- property_name`string`required:property name
- property_url`string` required:property url
- property_data `string` required:property data
- number_properties`string`required: total issue amount
- currency_identifier_desired`string`required:
- Deadline`string`required:deadline
- earlybird_bonus`string`required:weekly early birds bonus percentage
- total_number`string`required:crowdsale count

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout`  past output index
  - `scriptPubKey` lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","ecosystem":"1","precision":"5","previous_property_id":"","property_category":"wormhole wallet api","property_subcategory":"whcwallet","property_name":"whcwallet","property_url":"https://cash.wormhole.org","property_data":"welcome to wormhole","number_properties":"12345.987654","currency_identifier_desired":"1","deadline":"1535453483","earlybird_bonus":"10","total_number":"1029932"}' http://127.0.0.1:9999/getunsigned/51

//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff023fcfa435000000001976a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac0000000000000000896a4c86087768630000003301000600000000776f726d686f6c652077616c6c6574206170690077686377616c6c65742d676f0077686377616c6c65740068747470733a2f2f636173682e776f726d686f6c652e6f72670077656c636f6d6520746f20776f726d686f6c6500000000010000011f73d22358000000005b85292b0a00000000efccbb230000000000"}}
```

##### `/getunsigned/53`

Return the unsigned transaction of close crowdsale ,sign_data and unsigned_tx will be returned.  

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/53 |

**Argument**

- transaction_version`string`required: transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required: transaction issuer address

- currency_identifier`string`required:perperty id

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout` past output index
  - `scriptPubKey`  lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","currency_identifier":"23"}' http://127.0.0.1:9999/getunsigned/53
//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff020dd4a435000000001976a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac00000000000000000e6a0c08776863000000350000001700000000"}}
```

##### `/getunsigned/54`

Return the  unsigned transaction of creating fixed managed property ,sign_data and unsigned_tx will be returned. 

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/54 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:transaction issuer address
- ecosystem`string`required:location system
- precision`string` required:precision
- previous_property_id `string`required: previous property id
- property_category `string` required:property category
- property_subcategory`string`required:property subcategory
- property_name`string`required:property name
- property_url`string` required:property url
- property_data `string` required:property data

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout`  past output index
  - `scriptPubKey`  lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","ecosystem":"1","precision":"5","previous_property_id":"","property_category":"wormhole wallet api","property_subcategory":"whcwallet","property_name":"whcwallet","property_url":"https://cash.wormhole.org","property_data":"welcome to wormhole",}' http://127.0.0.1:9999/getunsigned/54

//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff02f7d0a435000000001976a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac00000000000000005d6a4c5a08776863000000360100080000000068656c6c6f20776f726c6400676f6c616e670077686377616c6c65740068747470733a2f2f636173682e776f726d686f6c652e6f72670077656c636f6d6520746f20776f726d686f6c650000000000"}}
```

##### `/getunsigned/55`

Return unsigned transaction of increasing circulation ,sign_data and unsigned_tx will be returned, Can only be used to manage property.

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/55 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:transaction issuer address 

- currency_identifier`string`required:porperty id
- Amount`string`required:increase count
- note`string`:remark

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout`  past output index
  - `scriptPubKey`  lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","currency_identifier":"3","amount":"121231.90","note":"hello world"}' http://127.0.0.1:9999/getunsigned/55

//Result
{"result":{"sign_data":[{"txid":"440580de6a14145bc151ef6b9811c374f4b1e298aa2d5956eff44898890f439f","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":0.00000546},{"txid":"9f40f5cdf0766dc4dcbcfeda2d25cff5103d96beab708497d556cf07e7fca688","vout":0,"scriptPubKey":"76a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac","value":8.99997721}],"unsigned_tx":"02000000029f430f899848f4ef56592daa98e2b1f474c311986bef51c15b14146ade8005440000000000ffffffff88a6fce707cf56d5978470abbe963d10f5cf252ddafebcdcc46d76f0cdf5409f0000000000ffffffff0245d3a435000000001976a9141fdcfa4003067f2326a92d1e582164f6f7792b0088ac0000000000000000226a20087768630000003700000003000000000001d98f68656c6c6f20776f726c640000000000"}}
```

##### `/getunsigned/56`

Return unsigned transaction of revoke property ,sign_data and unsigned_tx will be returned, Can only be used to manage property.

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/56 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key 
- fee`string` required:transaction fee
- transaction_from`string` required:transaction issuer address 

- currency_identifier`string`required:property id
- amount`string`required:revoke count
- note`string`:remark

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout` past output index
  - `scriptPubKey`  lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","currency_identifier":"3","amount":"121231.90","note":"hello world"}' http://127.0.0.1:9999/getunsigned/56

```

##### `/getunsigned/68`

Return unsigned transaction of burning BCH ,sign_data and unsigned_tx will be returned. 

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/68 |

**Argument**

- transaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:transaction issuer address 
- redeem_address`string`:change address
- amount_for_burn`string` required:burning amount

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout` past output index
  - `scriptPubKey`  lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","redeem_address":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","amount_for_burn":"1"}' http://127.0.0.1:9999/getunsigned/68
//Result
{"code":0,"message":"","result":{"sign_data":[{"txid":"a141d3f0cdceff520c9f041d6eb9e280f3cd91e1f5961015fa6567844e2c2a9e","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"31f17298a070da426babc44ce08cde63a482301dc6c1824ec8100a1ba49d9ad7","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"70123e418b0dd1252dfcf14cebde83b99aaaebe1f3472dc088747b572377ee15","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"f706c3d18625cabf4d1f512f5bd91f8d68bfd30fa49e6f7a33e9bfaea7f1d1ec","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"f8ef2594d1c051c0b0140f6322877eeefe4159ddc2e28f439ecf0a807cdb9a4a","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"4210e7faad9110fe9ea8a062010051ed1825dbec5bbbedcfff7c11b58196354d","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"990bb45c227364c0122ec3fbdd28961617cfb178d7c24d02fb91847e42fd8afc","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"c7cc58ea518044fa8283d1c48fc7b103efdb6a6d4c90f7c46d51728842ab1f14","vout":0,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.08999279},{"txid":"165843ed83737aee8b9d4a5f176d947fd4aa62fd1853f19feb0caeacdb398087","vout":2,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":0.00000546},{"txid":"9ded61176f43c929aa31e18aa98c4ac4b16197113d2feebc37e78b08242d7658","vout":0,"scriptPubKey":"76a914fdace26851d24a1b4b31dc50aee973bebc9d35b488ac","value":6.81988672}],"unsigned_tx":"020000000a9e2a2c4e846765fa151096f5e191cdf380e2b96e1d049f0c52ffcecdf0d341a10200000000ffffffffd79a9da41b0a10c84e82c1c61d3082a463de8ce04cc4ab6b42da70a09872f1310200000000ffffffff15ee7723577b7488c02d47f3e1ebaa9ab983deeb4cf1fc2d25d10d8b413e12700200000000ffffffffecd1f1a7aebfe9337a6f9ea40fd3bf688d1fd95b2f511f4dbfca2586d1c306f70200000000ffffffff4a9adb7c800acf9e438fe2c2dd5941feee7e8722630f14b0c051c0d19425eff80200000000ffffffff4d359681b5117cffcfedbb5becdb2518ed51000162a0a89efe1091adfae710420200000000fffffffffc8afd427e8491fb024dc2d778b1cf17169628ddfbc32e12c06473225cb40b990200000000ffffffff141fab428872516dc4f7904c6d6adbef03b1c78fc4d18382fa448051ea58ccc70000000000ffffffff878039dbacae0ceb9ff15318fd62aad47f946d175f4a9d8bee7a7383ed4358160200000000ffffffff58762d24088be737bcee2f3d119761b1c44a8ca98ae131aa29c9436f1761ed9d0000000000ffffffff0300e1f505000000001976a9140000000000000000000000000000000000376e4388ac00000000000000000a6a080877686300000044bb6b3723000000001976a914e6886aa39d558ecfa79981db118f238f8ca2308488ac00000000"}}
```

##### `/getunsigned/70`

Return the unsigned transaction of change issuer ,sign_data and unsigned_tx will be returned. 

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /getunsigned/70 |

**Argument**

- ransaction_version`string`required:transaction version number
- pubkey`string`:public key
- fee`string` required:transaction fee
- transaction_from`string` required:transaction issuer address 
- transaction_to`string` required:transaction receive address
- currency_identifier`string` required:property id

**Result**

- `unsigned_tx`：serialization of unsigned transactions
- `sign_data`
  - `txid`  transaction hash
  - `vout`  past output index
  - `scriptPubKey`  lock script
  - `value`  value

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"transaction_version":"0","pubkey":"","fee":"0.0001","transaction_from":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","transaction_to":"bchtest:qpncl685d5eys3tj4vykrhw4507fcuas9scm4qd70e","currency_identifier":"203"}' http://127.0.0.1:9999/getunsigned/70

//Result
{"result":{"sign_data":[{"txid":"522d16cfe9392eb43d58ddb4b450505e5ad4ded2a2d6ac4f0526c81665e4ad4a","vout":2,"scriptPubKey":"76a9140ee020c07f39526ac5505c54fa1ab98490979b8388ac","value":0.00000546},{"txid":"fd998ed2567d77428a85687b0ccc278fffe9a813d02a98d0cfccc3d42811a282","vout":2,"scriptPubKey":"76a9140ee020c07f39526ac5505c54fa1ab98490979b8388ac","value":0.00000546},{"txid":"3c17c8e82ff8c5b1eb7109b0c554dfbda4dd7a9b69cc67f7917c77f6eab580c0","vout":0,"scriptPubKey":"76a9140ee020c07f39526ac5505c54fa1ab98490979b8388ac","value":0.0001}],"unsigned_tx":"02000000034aade46516c826054facd6a2d2ded45a5e5050b4b4dd583db42e39e9cf162d520200000000ffffffff82a21128d4c3cccfd0982ad013a8e9ff8f27cc0c7b68858a42777d56d28e99fd0200000000ffffffffc080b5eaf6777c91f767cc699b7adda4bddf54c5b00971ebb1c5f82fe8c8173c0000000000ffffffff03e8130000000000001976a9140ee020c07f39526ac5505c54fa1ab98490979b8388ac00000000000000000e6a0c0877686300000046000000cb22020000000000001976a914678fe8f46d32484572ab0961ddd5a3fc9c73b02c88ac00000000"}}
```

##### `/transaction/push`

To broadcasting the signed wormhole transactions.

| Protocol | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | POST   | /getunsigned/4 |

**Argument**

- signedTx`string` required:hash value of the signed transactions

**Result**

- `txHash`:return the hash value  after broadcasting

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"tx":"02000000014f01e00288383dd8d432c716802e20ba67792ed2fa66968827b522b5f94aac49000000006a47304402204e30a6321becdabe61c8985d719e94f1ba978bc58cef3b6d7bdd1794e2143e6c02204821be57779140003156f07cd6c8698719a8015df8a3a137f611c9fb3b71422b412103ac7fa295271559c620d0d24fa6754f14d33768d1af959ed5fd1ada4ce117315bffffffff03fe8a8d06000000001976a9140ee020c07f39526ac5505c54fa1ab98490979b8388ac0000000000000000526a4c4f087768630000003601000100000000436f6d70616e69657300426974636f696e2043617368204d696e696e6700514d43007777772e716d632e63617368004d616465207769746820424954424f580022020000000000001976a9140ee020c07f39526ac5505c54fa1ab98490979b8388ac00000000"}' http://127.0.0.1:9999/transaction/push
//Result
{ "result": { "txHash": "ebfe99a35ab72fb0dbee7bcecd61ac26e636a52984bd2224748856145b737d84"
 }}
```

##### `/env`

Return to the current network environment.

| Protocol | Method | API  |
| -------- | ------ | ---- |
| HTTP     | GET    | /env |

**Argument**

- None

**Result**

- `testnet`:network environment

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/env

//Result
{"result":{"testnet":true}}
```

##### `/category`

Return the default classification list for the back end.

| Protocol | Method | API       |
| -------- | ------ | --------- |
| HTTP     | GET    | /category |

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/category
//Result
{"result":{"Accommodation and food":["Accommodation","Food and beverage","Other"],"Administrative and support":["Rental and leasing activities","Employment activities","Security and investigation activities","Other"],"Agriculture, forestry and fishing":["Forestry and logging","Fishing and aquaculture","Other"],"Arts, entertainment and recreation":["Libraries, archives, museums","Gambling and betting activities","Sports activities and amusement","Other"],"Construction":["Construction of buildings","Civil engineering","Specialized construction","Other"],"Education":["Education","Other"],"Extraterritorial organizations":["Extraterritorial organizations","Other"],"Information and communication":["Publishing activities","Programming and broadcasting activities","Telecommunications","Information service activities","Other"],"Manufacturing":["food products","beverages","tobacco products","textiles","wearing apparel","leather and related products","paper products","reproduction of recorded media","chemical products","rubber and plastics products","basic metals","computer, electronic","electrical equipment","machinery and equipment n.e.c.","other transport equipment","furniture","Other manufacturing","Other"],"Mining and quarrying":["Forestry and logging","Mining of metal ores","Other mining and quarrying","Mining support service","Other"],"Other":["Other"],"Real estate activities":["Real estate activities","Other"],"Transportation and storage":["Water transport","Air transport","Postal and courier","Other"],"social work activities":["Human health activities","Residential care activities","Other"]}}
```

##### `/category/subcategories?category=***`

Return a list of subcategories under a main category

| Protocol | Method | API                                  |
| -------- | ------ | ------------------------------------ |
| HTTP     | GET    | /category/subcategories?category=*** |

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/category/subcategories?category=social%20work%20activities

//Request
{"result":["Human health activities","Residential care activities","Other"]}
```

##### `/balance/addresses`

Return the property balance information corresponding to the address list

| Protocol | Method | API                |
| -------- | ------ | ------------------ |
| HTTP     | POST   | /balance/addresses |

**Argument**

- address`string` required:address to be queried

**Result**

- `address`current address
- `property_id`property id
- `property_name`property name
- `precision`precision
- `balance_available`balance available

**Example**

```go 
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg"}' http://127.0.0.1:9999/balance/addresses

//Result
{"result":{"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg":[{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":1,"property_name":"WHC","precision":8,"balance_available":"18.6688386","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":108,"property_name":"My Very First Token","precision":1,"balance_available":"1496","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":109,"property_name":"Quantum Miner","precision":1,"balance_available":"10000","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":110,"property_name":"Quantum Miner","precision":1,"balance_available":"0","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":111,"property_name":"Quantum Miner","precision":1,"balance_available":"999000","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":112,"property_name":"Time Machine","precision":1,"balance_available":"0","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":115,"property_name":"luzhiyao","precision":3,"balance_available":"6.169","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":152,"property_name":"TTK","precision":1,"balance_available":"1000000","pendingpos":"0","pendingneg":"0"},{"address":"bchtest:qq2j9gp97gm9a6lwvhxc4zu28qvqm0x4j5e72v7ejg","property_id":180,"property_name":"TTK","precision":0,"balance_available":"100","pendingpos":"0","pendingneg":"0"}]}}
```

##### `/history/list`

Return the transaction record corresponding to the address list.

| Protocol | Method | API           |
| -------- | ------ | ------------- |
| HTTP     | POST   | /history/list |

**Argument**

- address`array` required:address to be queried
- property_id`string`required:

**Result**

- `list`
  - `tx_hash`：transaction hash
  - `tx_type`：transaction type
  - `tx_state`transaction state
  - `address`address to be queried
  - `address_role`address role
  - `balance_available_credit_debit`balance
  - `property_id`property id
  - `property_name`property name
  - `created`transaction creation time

- `total`eligible results amount

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"address":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","property_id":"345"}' http://127.0.0.1:9999/history/list
//Result
{{"result":{"list":[{"tx_hash":"d838466af4815c5a6ea131ad3ad3ec0ff8cf9a850d1045a61e609d72840db531","tx_type":3,"tx_state":"valid","address":"bchtest:qq0ae7jqqvr87gex4yk3ukppvnm0w7ftqqpzv0lcqa","address_role":"payee","balance_available_credit_debit":"0.03595587","property_id":1,"property_name":"WHC","created":1544619708}]},"total":1}}
```

##### `/history/id/:id`

Return the transaction history of a property.

| Protocol | Method | API             |
| -------- | ------ | --------------- |
| HTTP     | POST   | /history/id/:id |

**Argument**

- Start`string`required:start index
- Count`string`required:count of records
- Address`string`required:address to be queried

**Result**

- `total`eligible results amount

- `transactions`
  - `amount`：transaction amount
  - `block`blocks of the current transaction
  - `blockhash`block hash
  - `blocktime`block time
  - `confirmations`confirmation times
  - `fee`transaction fee
  - `fee_rate`fee rate
  - `precision`precision
  - `propertyid`  property id
  - `referenceaddress`transaction receiver address
  - `sendingaddress` transaction issuer address
  - `txid` transaction hash
  - `type`transaction type
  - `type_int`transaction type number
  - `valid`effectiveness
  - `version`version

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"address":"bchtest:qrv2054hm0rswdnyaskxsyxaq6j58x479ynsgl79xp","start":"0"，"count":"5"}' http://127.0.0.1:9999/history/id/1
//Result
{"result":{"total":0,"transactions":[{"amount":"1","block":1262852,"blockhash":"00000000000006631dd79dfed47761a9823b8e9efe977361e6de3f2cdda0c258","blocktime":1539678705,"confirmations":8670,"fee":"364","fee_rate":"0.00001421","ismine":false,"precision":"8","propertyid":345,"referenceaddress":"bchtest:qqjpd98pj4h464yjhwe6svsq6c5qu6h2rsts63p7un","sendingaddress":"bchtest:qrv2054hm0rswdnyaskxsyxaq6j58x479ynsgl79xp","subsends":null,"txid":"39e5676d313fcf9c6a85eddf7b69bee67b22adc429e193d1978cff20b912ac35","type":"Simple Send","type_int":0,"valid":true,"version":0}}
```

##### `/history/detail?tx_hash=***`

Return the details of a transaction

| Protocol | Method | API                         |
| -------- | ------ | --------------------------- |
| HTTP     | GET    | /history/detail?tx_hash=*** |

****

**Argument**

- tx_hash`string` required:transaction hash value
- type`string` required:transaction type

**Result**

- `result`
  - `hex`serialization of  transactions
  - `txid`transaction hash
  - `version`transaction version
  - `locktime`delay Time
  - `vin`
    - `txid`past transaction hash
    - `vout`past output index
    - `scriptSig`Unlock script
    - `sequence`serial number
  - `vout`
    - `value`transaction amount
    - `n` index 
    - `scriptPubKey`lock script
  - `blockhash`block hash
  - `confirmations`confirmation times
  - `time  `transaction time
  - `blocktime` block time 

**Example**

```go 
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/history/detail?tx_hash=04b1c47a23ae63e9931576b921427fc0f4ad26fa29c4606cdc64a1efe6636cd9&type=bch
//Result
{"result":{"hex":"02000000018327d3cabd19ee42c9475c15c3516eef2f758fbbdf7953bea2011868025a40e8010000006b4830450221008371cfc96857084c95b52c285b50f9c3544d1e2e21f433a992873c3f87be244002203153354ba88f14f5b89aeb9fbd7aa06dfd3a72a873cc5a9e867ef6c97a61f44a4121039aece7c9a048d922468185b2b51c4265443f08d2515b27d0d1f73abe1371e171ffffffff034e4e8e06000000001976a914d8a7d2b7dbc7073664ec2c6810dd06a5439abe2988ac0000000000000000166a1408776863000000000000000100000000069f406022020000000000001976a914c23c0b8692c1a27721ea2ec787e7eab553b0d07088ac00000000","txid":"04b1c47a23ae63e9931576b921427fc0f4ad26fa29c4606cdc64a1efe6636cd9","version":2,"locktime":0,"vin":[{"txid":"e8405a02681801a2be5379dfbb8f752fef6e51c3155c47c942ee19bdcad32783","vout":1,"scriptSig":{"asm":"30450221008371cfc96857084c95b52c285b50f9c3544d1e2e21f433a992873c3f87be244002203153354ba88f14f5b89aeb9fbd7aa06dfd3a72a873cc5a9e867ef6c97a61f44a[ALL|FORKID] 039aece7c9a048d922468185b2b51c4265443f08d2515b27d0d1f73abe1371e171","hex":"4830450221008371cfc96857084c95b52c285b50f9c3544d1e2e21f433a992873c3f87be244002203153354ba88f14f5b89aeb9fbd7aa06dfd3a72a873cc5a9e867ef6c97a61f44a4121039aece7c9a048d922468185b2b51c4265443f08d2515b27d0d1f73abe1371e171"},"sequence":4294967295}],"vout":[{"value":1.09989454,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 d8a7d2b7dbc7073664ec2c6810dd06a5439abe29 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914d8a7d2b7dbc7073664ec2c6810dd06a5439abe2988ac","reqSigs":1,"type":"pubkeyhash","addresses":["bchtest:qrv2054hm0rswdnyaskxsyxaq6j58x479ynsgl79xp"]}},{"value":0,"n":1,"scriptPubKey":{"asm":"OP_RETURN 08776863000000000000000100000000069f4060","hex":"6a1408776863000000000000000100000000069f4060","type":"nulldata"}},{"value":0.00000546,"n":2,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 c23c0b8692c1a27721ea2ec787e7eab553b0d070 OP_EQUALVERIFY OP_CHECKSIG","hex":"76a914c23c0b8692c1a27721ea2ec787e7eab553b0d07088ac","reqSigs":1,"type":"pubkeyhash","addresses":["bchtest:qrprczuxjtq6yaepaghv0pl8a2648vxswqxzdlnety"]}}],"blockhash":"00000000000a8cbe9e0c1d3f36a8a68f80889618f5c834552b7db8c096680aa1","confirmations":11636,"time":1539347623,"blocktime":1539347623}}
```

##### `/history/detail/pending`

Return details of pending wormhole transactions.

| Protocol | Method | API                     |
| -------- | ------ | ----------------------- |
| HTTP     | GET    | /history/detail/pending |

**Argument**

- tx_hash`string` required:transaction hash value

**Result**

- `txid` transaction hash value
- `fee`transaction fee
- `sendingaddress`transaction  issuer address
- `referenceaddress`transaction receiver address
- `version`version
- `type_int`transaction type 
- `type`transaction type number
- `block`blocks of the Current transaction
- `confirmations`confirmation time
- `valid`effectiveness
- `blockhash`block hash
- `blocktime `block time
- `propertyid` property id
- `precision`precision
- `amount` transaction amount

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/history/detail/pending?tx_hash=bc06ccd0f9d401f674558545bb79e1c20b61807184226578860641c87422430a

//Result
{"result":{"txid":"bc06ccd0f9d401f674558545bb79e1c20b61807184226578860641c87422430a","fee":"10000","sendingaddress":"bchtest:qrv2054hm0rswdnyaskxsyxaq6j58x479ynsgl79xp","referenceaddress":"bchtest:qrprczuxjtq6yaepaghv0pl8a2648vxswqxzdlnety","ismine":false,"version":0,"type_int":0,"type":"Simple Send","block":1262062,"confirmations":11758,"valid":true,"blockhash":"00000000003a7e463e04898c464717f9b10a1dfe864b6b63044bf395f7c9f60e","blocktime":1539264034,"propertyid":345,"precision":"8","amount":"1.09879800","subsends":null}}
```

##### `/property/id/:id`

| Protocol | Method | API              |
| -------- | ------ | ---------------- |
| HTTP     | GET    | /property/id/:id |

**Argument**

- None

**Result**

- `category` property category
- `creationtxid `property creation hash
- `data` remark
- `fixedissuance`fixed issuance
- `freezingenabled` freezing enabled
- `issuer`property creator address
- `managedissuance` managed issuance
- `name`property name
- `precision` precision
- `propertyid` property id
- `subcategory` sub category
- `totaltokens` total token quantity
- `url` url

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/property/id/1
//Result
{"code":0,"message":"","result":{"category":"Information and communication","creationtxid":"5f03e20f1830aa68266ef401ad053bd2592df58bc2f944015e7dcfe4fefc632e","data":"wormhole and bitcoin cash will become better","fixedissuance":true,"freezingenabled":null,"issuer":"bchtest:qrv2054hm0rswdnyaskxsyxaq6j58x479ynsgl79xp","managedissuance":false,"name":"whcwallet-test","precision":8,"propertyid":345,"subcategory":"Telecommunications","totaltokens":"10000.12345678","url":"https://github.com/bcext"}}
```

##### `/property/listowners/:id`

Return to the owner of the property list. 

| Protocol | Method | API                      |
| -------- | ------ | ------------------------ |
| HTTP     | POST   | /property/listowners/:id |

**Argumment**

- pageNo `string` required:page number(can be null)
- pageSize `string` required:page size(can be null)

**Result**

- `Address`the address of the owner of the property
- `balance_available`balance available
- `Status`confirmation status

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST  http://127.0.0.1:9999/property/listowners/459
//Result
{"result":{"list":[{"Address":"bchtest:qq7pl7ta6uf2kpjphf3p22l0jqxgeyrx9vjx29dr8t","balance_available":"9800","Status":0},{"Address":"bchtest:qpz8py2yqyp7x2aaqfvrwdk4jf2c23ypzczr6weclj","balance_available":"100","Status":0},{"Address":"bchtest:qqu9lh4jpc05p59pfhu9amyv9uvder8j3sa2up95vs","balance_available":"100","Status":1}],"total":3}}
```

##### `/property/name`

Returns the property corresponding to the name.

| Property | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | GET    | /property/name |

**Argument**

- None

**Result**

- property detail

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/property/name/zc_crow2
//Result
{"result":[{"active":true,"addedissuertokens":"0.0000000","amountraised":"64.77041915","category":"","closetx":null,"creationtxid":"5cba679e754f038062da9becf7f2f41d5cd093cde8bd306179611d27a02f08cd","data":"hi","deadline":1596127632,"earlybonus":10,"fixedissuance":false,"freezingenabled":null,"issuer":"bchtest:qz3fgledq0tgl0ry6pn0c5nmufspfrr8aqsyuc39yl","managedissuance":false,"name":"zc_crow2","participanttransactions":null,"precision":"7","propertyid":27,"propertyiddesired":1,"purchasedtokens":665.9166694,"starttime":1532970004,"subcategory":"","tokensissued":"100000","tokensperunit":"1","totaltokens":"100000","url":"www.zc"}]}
```

##### `/property/query/:keyword`

Return the property details ，support the fuzzy search via property name or property id

| Protocol | Method | API                      |
| -------- | ------ | ------------------------ |
| HTTP     | POST   | /property/query/:keyword |

**Argument**

- keyword`string` required:key word (can be null)

**Result**

- property detail

**Example**

```go
//Request
curl -i -H 'Content-Type: application/json' -X POST -d{"keyword":"test"} http://127.0.0.1:9999/property/query/
//Result
{"result":[{"category":"aaaatest","creationtxid":"5e57721389707e2c8594518628ba9487a74a89e250e8b21df820e327d7b8eefd","data":"Test Token for AAAA","fixedissuance":false,"freezingenabled":false,"issuances":[],"issuer":"bchtest:qpv62pqneh59wvm59c8p3sxu9u00n86psysnr9a909","managedissuance":true,"name":"aaaa test token","precision":1,"propertyid":65,"subcategory":"aaaa test","totaltokens":"0","url":"www.aaaatest.com"}]}
```

##### `/property/address/:address`

Returns the property of an address.

| Property | Method | API                        |
| -------- | ------ | -------------------------- |
| HTTP     | GET    | /property/address/:address |

**Argument**

- None

**Result**

- `balance`

- `property`
  - `category` property category
  - `creationtxid `property creation hash
  - `data` remark
  - `fixedissuance`fixed issuance
  - `freezingenabled` freezing enabled
  - `issuer`property creator address
  - `managedissuance` managed issuance
  - `name`property name
  - `precision` precision
  - `propertyid` property id
  - `subcategory` sub category
  - `totaltokens` total token quantity
  - `url` url

**Example**

```go 
//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/property/address/bchtest:qz3fgledq0tgl0ry6pn0c5nmufspfrr8aqsyuc39yl
//Result
{"result":[{"balance":7.764596525e-7,"property":{"category":"N/A","creationtxid":"0000000000000000000000000000000000000000000000000000000000000000","data":"WHC serve as the binding between Bitcoin cash, smart properties and contracts created on the Wormhole.","fixedissuance":false,"freezingenabled":null,"issuer":"bchtest:qqqqqqqqqqqqqqqqqqqqqqqqqqqqqdmwgvnjkt8whc","managedissuance":false,"name":"WHC","precision":8,"propertyid":1,"subcategory":"N/A","totaltokens":"47613.45013","url":"http://www.wormhole.cash"}}}
```

##### `/property/list`

Returns all proerty of the current system.support fuzzy search via property name or property id

| Property | Method | API            |
| -------- | ------ | -------------- |
| HTTP     | GET    | /property/list |

**Argument**

- pageNo `string` required:page number(can be null)
- pageSize `string` required:page size(can be null)
- Keyword`string`required:key word

**Result**

- `list`
  - `category` property category
  - `creationtxid `property creation hash
  - `data` remark
  - `fixedissuance`fixed issuance
  - `freezingenabled` freezing enabled
  - `issuer`property creator address
  - `managedissuance` managed issuance
  - `name`property name
  - `precision` precision
  - `propertyid` property id
  - `subcategory` sub category
  - `totaltokens` total  token quantity
  - `url` url

**Example**

```go

//Request
curl -i -H 'Content-Type: application/json' -X GET http://127.0.0.1:9999/property/list?pageSize=5&pageNo=1&keyword=22
//Result
{"result":{"list":[{"category":"test","created":1532937514,"creationtxid":"e2c320795cdf8c9ba3c9404a22c75b444f373b433f43980c7a5e2bac3a994a4e","data":"mydata2","fixedissuance":true,"freezingenabled":null,"issuer":"bchtest:qrzrgmqfdn42j6vj8l47a6ckzgw2pxt7xud89uwgfn","managedissuance":false,"name":"snow2","precision":1,"propertyid":10,"subcategory":"test2","totaltokens":"2100000000000000","url":"www.snow2.token"}]}}
```

`/property/listbyowner`

Return the property whose issuer is the specified address

| Property | Method | API                   |
| -------- | ------ | --------------------- |
| HTTP     | POST   | /property/listbyowner |

**Argument**

- pageNo `string` required:page number(can be null)
- pageSize `string` required:page size(can be null)
- address`string` required:address to be queried

**Result**

- `list`
  - `TxData`
    - `amount`transaction amount
    - `block`blocks of the current transaction
    - `blockhash`block hash
    - `blocktime`block time
    - `confirmations` confirmations count
    - `data`remark
    - `ecosystem` system
    - `fee`transaction fee
    - `fee_rate` fee rate
    - `precision`precision
    - `propertyid`property id
    - `propertyname`property name
    - `sendingaddress` transaction issuer address
    - `totalstofee`total fee
    - `txid`transaction hash
    - `type`transaction type
    - `type_int`transaction type number
    - `url`url 
    - `valid`effectiveness
    - `version`version
  - `PropertyData`
    - `category`category
    - `creationtxid`create transaction hash
    - `data`remark
    - `fixedissuance` property type
    - `freezingenabled` freezing enabled
    - `issuer`issuer
    - `managedissuance` property type
    - `name`property name
    - `precision`precision
    - `propertyid` property id
    - `subcategory`subcategory
    - `totaltokens` total tokens
    - `url`url
- `total`

**Example**

```go
//Request
curl -i -H 'Content-Type: application/x-www-form-urlencoded' -X POST -d'{"address":"bchtest:qrv2054hm0rswdnyaskxsyxaq6j58x479ynsgl79xp","pageSize":"5"，"pageNo":"1"}' http://127.0.0.1:9999/property/listbyowner

//Result
{,"result":{"list":[{"TxData":{"amount":"1.0","block":1261784,"blockhash":"0000000000000a076f5dfacf25e5df108b65a340b4a3ffc49ea1bb4c310718d0","blocktime":1539153861,"confirmations":9738,"data":"1","ecosystem":"main","fee":"2279","fee_rate":"0.000086","ismine":false,"precision":"1","propertyid":337,"propertyname":"1","sendingaddress":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","subsends":null,"totalstofee":"1","txid":"986c0e274b1bb1d9167fa15bb3bec616f558e6e50b63ffeeee08bf8f25663ce2","type":"Create Property - Fixed","type_int":50,"url":"1","valid":true,"version":0},"PropertyData":{"category":"","creationtxid":"986c0e274b1bb1d9167fa15bb3bec616f558e6e50b63ffeeee08bf8f25663ce2","data":"1","fixedissuance":true,"freezingenabled":null,"issuer":"bchtest:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk","managedissuance":false,"name":"1","precision":1,"propertyid":337,"subcategory":"","totaltokens":"1","url":"1"}}]}
```

##### `headers/file`

Download the blockheader file for the main network.

| Protocol | Method | API          |
| -------- | ------ | ------------ |
| HTTP     | GET    | headers/file |
