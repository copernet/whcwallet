package model

import (
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// AddressRole type definitions:
// sender sent a Simple Send
// recipient received a Simple Send
// issuer created a smart property
// participant is an investor in a crowdsale
// payee received a STO amount
// seller created a DEx sell offer
// buyer accepted a DEx sell offer

type AddressRole string

// Notice:
// Crowdsale Purchase will create two transactions history:
// 1. AddressRole: sender
// 2. AddressRole: participant
//
// Create property will create two transaction history:
// 1. AddressRole: feepayer
// 2. AddressRole: issuer
const (
	Sender      AddressRole = "sender"      // Crowdsale Purchase || sender （Simple Send transaction）
	Recipient   AddressRole = "recipient"   // receiver(tx outs) （Simple Send transaction）
	Issuer      AddressRole = "issuer"      // create Property
	Participant AddressRole = "participant" // Crowdsale Purchase
	Payee       AddressRole = "payee"       // receive property from (send to owner transaction)
	Buyer       AddressRole = "buyer"       // Burn BCH Get WHC
	Feepayer    AddressRole = "feepayer"    // issue (Create Property - Fixed)
	Payer       AddressRole = "payer"       // issue Send To Owners transaction
	Seller      AddressRole = "seller"
)

func (a *AddressRole) Scan(value interface{}) error {
	*a = AddressRole(value.([]byte))
	return nil
}

func (a AddressRole) Value() (driver.Value, error) {
	return string(a), nil
}

type Protocol string

const (
	Fiat        Protocol = "fiat"
	BitcoinCash Protocol = "bitcoincash"
	Wormhole    Protocol = "wormhole"
)

func (p *Protocol) Scan(value interface{}) error {
	*p = Protocol(value.([]byte))
	return nil
}

func (p Protocol) Value() (driver.Value, error) {
	return string(p), nil
}

type Ecosystem string

const (
	Production Ecosystem = "production"
	TestSystem Ecosystem = "testsystem"
)

func (e *Ecosystem) Scan(value interface{}) error {
	*e = Ecosystem(value.([]byte))
	return nil
}

func (e Ecosystem) Value() (driver.Value, error) {
	return string(e), nil
}

type TxState string

const (
	Pending TxState = "pending"
	Valid   TxState = "valid"
	InValid TxState = "invalid"
)

func (t *TxState) Scan(value interface{}) error {
	*t = TxState(value.([]byte))
	return nil
}

func (t TxState) Value() (driver.Value, error) {
	return string(t), nil
}

type Client string

const (
	Android Client = "android"
	IOS     Client = "ios"
)

func (t *Client) Scan(value interface{}) error {
	*t = Client(value.([]byte))
	return nil
}

func (t Client) Value() (driver.Value, error) {
	return string(t), nil
}

//Database Table Structure

type Block struct {
	ID          uint      `json:"-" gorm:"primary_key"`
	Version     int32     `json:"version"`
	BlockHeight int64     `json:"block_height" gorm:"not null;index"`
	BlockHash   string    `json:"block_hash" gorm:"type:char(64);index"`
	Nonce       uint32    `json:"nonce"`
	Bits        uint32    `json:"bits"`
	PrevBlock   string    `json:"prev_block" gorm:"type:char(64)"`
	MerkleRoot  string    `json:"merkleroot"`
	BlockTime   int64     `json:"block_time"`
	Txcount     int       `json:"txcount"`
	Whccount    int       `json:"whccount"`
	Size        int32     `json:"size"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// all the transactions we know about; keeping them even after an address or wallet is de-activated
// only wormhole transaction

type Tx struct {
	ID            uint      `json:"-" gorm:"primary_key"`
	TxID          int       `json:"tx_id" gorm:"not null;index"`
	TxHash        string    `json:"tx_hash" gorm:"type:char(64);not null;index"`
	Protocol      Protocol  `json:"protocol" gorm:"not null;default:'bitcoincash'" sql:"type:ENUM('fiat', 'bitcoincash', 'wormhole')"`
	TxType        uint64    `json:"tx_type" gorm:"not null"`
	Ecosystem     Ecosystem `json:"ecosystem" sql:"type:ENUM('production', 'testsystem')" gorm:"default:'production'"`
	TxState       TxState   `json:"tx_state" gorm:"not null;default:'pending'" sql:"type:ENUM('pending', 'valid', 'invalid')"`
	TxErrorCode   int16     `json:"tx_error_code" gorm:"default: 0"`
	TxBlockHeight uint32    `json:"tx_block_height" gorm:"default: 0"`
	TxSeqInBlock  int       `json:"tx_seq_in_block" gorm:"default: 0"`
	BlockTime     int64     `json:"block_time" gorm:"not null;index"`
	CreatedAt     time.Time `json:"-"`
}

// data that is specific to the particular transaction type, as a JSON object */
type TxJson struct {
	ID        uint      `json:"-" gorm:"primary_key"`
	TxID      int       `json:"tx_id" gorm:"not null;index"`
	Protocol  Protocol  `json:"protocol" gorm:"not null;default:'bitcoincash'" sql:"type:ENUM('fiat', 'bitcoincash', 'wormhole')"`
	TxData    string    `json:"tx_data" gorm:"type:json;not null"`
	RawData   string    `json:"raw_data" gorm:"type:text"`
	CreatedAt time.Time `json:"-"`
}

// Addresses with private keys owned by each Wallet. See Following table for objects watched by a wallet. */
type AddressesInTx struct {
	ID                          uint             `json:"-" gorm:"primary_key"`
	Address                     string           `json:"address" gorm:"type:varchar(64);not null;index"`
	PropertyID                  int64            `json:"property_id" gorm:"not null"`
	Protocol                    Protocol         `json:"protocol" gorm:"not null;default:'bitcoincash'" sql:"type:ENUM('fiat', 'bitcoincash', 'wormhole')"`
	TxID                        int              `json:"tx_id" gorm:"not null"`
	AddressTxIndex              int16            `json:"address_tx_index" gorm:"not null"`
	AddressRole                 AddressRole      `json:"address_role" gorm:"not null" sql:"type:ENUM('sender', 'recipient', 'issuer', 'participant', 'payee', 'seller', 'buyer','feepayer','payer')"`
	BalanceAvailableCreditDebit *decimal.Decimal `json:"balance_available_credit_debit" sql:"type:decimal(64,8)"`
	BalanceReservedCreditDebit  *decimal.Decimal `json:"balance_reserved_credit_debit" sql:"type:decimal(64,8)"`
	BalanceAcceptedCreditDebit  *decimal.Decimal `json:"balance_accepted_credit_debit" sql:"type:decimal(64,8)"`
	BalanceFrozenCreditDebit    *decimal.Decimal `json:"balance_frozen_credit_debit" sql:"type:decimal(64,8)"`
	CreatedAt                   time.Time        `json:"-"`
	UpdatedAt                   time.Time        `json:"-"`
}

//Balances for each PropertyID (currency) owned by an Address
//for all addresses we know about, even if they're not in a wallet
type AddressBalance struct {
	ID               uint             `json:"-" gorm:"primary_key"`
	Address          string           `json:"address" gorm:"type:varchar(64);index"`
	PropertyID       int64            `json:"property_id" gorm:"not null;default:0"`
	Protocol         Protocol         `json:"protocol" gorm:"not null;default:'bitcoincash'" sql:"type:ENUM('fiat', 'bitcoincash', 'wormhole')"`
	Ecosystem        Ecosystem        `json:"ecosystem" sql:"type:ENUM('production', 'testsystem')" gorm:"default:'production'"`
	BalanceAvailable *decimal.Decimal `json:"balance_available" gorm:"not null;default:0"  sql:"type:decimal(64,8)"`
	BalanceReserved  *decimal.Decimal `json:"balance_reserved" gorm:"not null;default:0"  sql:"type:decimal(64,8)"`
	BalanceAccepted  *decimal.Decimal `json:"balance_accepted" gorm:"not null;default:0"  sql:"type:decimal(64,8)"`
	BalanceFrozen    *decimal.Decimal `json:"balance_frozen" gorm:"not null;default:0"  sql:"type:decimal(64,8)"`
	LastTxID         int              `json:"last_tx_id"`
	CreatedAt        time.Time        `json:"-"`
	UpdatedAt        time.Time        `json:"-"`
}

// current state of Smart Properties; 1 row for each SP */
type SmartProperty struct {
	ID                  uint      `json:"-" gorm:"primary_key"`
	Protocol            Protocol  `json:"protocol" gorm:"not null;default:'bitcoincash'" sql:"type:ENUM('fiat', 'bitcoincash', 'wormhole')"`
	PropertyID          int64     `json:"property_id" gorm:"index"`
	Issuer              string    `json:"issuer" gorm:"type:varchar(64);not null"`
	Ecosystem           Ecosystem `json:"ecosystem" sql:"type:ENUM('production', 'testsystem')" gorm:"default:'production'"`
	CreateTxID          int       `json:"create_tx_id" gorm:"not null"`
	LastTxID            int       `json:"last_tx_id" gorm:"not null"`
	PropertyName        string    `json:"property_name" gorm:"type:varchar(256)"`
	Precision           int       `json:"precision"`
	PrevPropertyID      int64     `json:"prev_property_id" gorm:"DEFAULT:0"`
	PropertyServiceURL  string    `json:"property_service_url" gorm:"type:varchar(256)"`
	PropertyCategory    string    `json:"property_category" gorm:"type:varchar(256)"`
	PropertySubcategory string    `json:"property_subcategory" gorm:"type:varchar(256)"`
	RegistrationData    string    `json:"registration_data" gorm:"type:varchar(5000)"`
	PropertyData        string    `json:"property_data" gorm:"type:json"`
	Flags               string    `json:"flags" gorm:"type:text"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// the list of transactions that affected each SP */
type PropertyHistory struct {
	ID         uint      `json:"-" gorm:"primary_key"`
	PropertyID int64     `json:"property_id" gorm:"not null"`
	TxID       int64     `json:"tx_id" gorm:"not null;index"`
	CreatedAt  time.Time `json:"-"`
}

type Session struct {
	ID         uint      `json:"-" gorm:"primary_key"`
	SessionID  string    `json:"session_id" gorm:"type:varchar(64);not null;index"`
	Challenge  *string   `json:"challenge" gorm:"type:varchar(64)"`
	PChallenge *string   `json:"p_challenge" gorm:"type:varchar(64)"`
	PubKey     string    `json:"pub_key" gorm:"type:text"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

// Wallets have addresses with private keys. Objects being watched are in the Following table */
type Wallet struct {
	ID         uint      `json:"-" gorm:"primary_key"`
	WalletID   string    `json:"wallet_id" gorm:"index"`
	LastLogin  time.Time `json:"last_login"`
	LastBackup time.Time `json:"last_backup"`
	WalletBlob string    `json:"wallet_blob" gorm:"type:text"`
	UserName   string    `json:"user_name" gorm:"type:varchar(32)"`
	Email      string    `json:"email" gorm:"type:varchar(64);unique_index"`
	Ip         string    `json:"ip" gorm:"type:varchar(64)"`
	Settings   string    `json:"settings" gorm:"type:text"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type Notification struct {
	ID        uint   `json:"-" gorm:"primary_key"`
	Address   string `json:"address" gorm:"index"`
	Timestamp int64  `json:"timestamp"`
	HistoryID uint   `json:"-"`
}

type Version struct {
	ID          uint      `json:"-" gorm:"primary_key"`
	Device      Client    `json:"device" sql:"type:ENUM('android', 'ios')"`
	Version     string    `json:"version"`
	VersionCode     int    `json:"version_code"`
	Download    string    `json:"download"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// base struct vo
type BalanceNotify struct {
	Address    string
	PropertyID int64
	TxID       int
}
