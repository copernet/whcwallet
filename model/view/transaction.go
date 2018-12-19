package view

import "github.com/shopspring/decimal"

// tx type: 0
type SimpleSendTx struct {
	TxVersion       string   `form:"transaction_version" binding:"required"`
	PubKey          *string  `form:"pubkey" binding:"-"`
	Fee             float64  `form:"fee" binding:"gte=0.00001,required"`
	From            string   `form:"transaction_from" binding:"required"`
	To              string   `form:"transaction_to" binding:"required"`
	CID             int64    `form:"currency_identifier" binding:"gte=1,required"`
	Amount          float64  `form:"amount_to_transfer" binding:"gt=0,required"`
	RedeemAddress   *string  `form:"redeem_address" binding:"-"`
	ReferenceAmount *float64 `form:"reference_amount"`
}

// tx type: 1
type ParticipateCrowdSale struct {
	TxVersion       string   `form:"transaction_version" binding:"required"`
	PubKey          *string  `form:"pubkey" binding:"-"`
	Fee             float64  `form:"fee" binding:"gte=0.00001,required"`
	From            string   `form:"transaction_from" binding:"required"`
	To              string   `form:"transaction_to" binding:"required"`
	Amount          float64  `form:"amount_to_transfer" binding:"gt=0,required"`
	RedeemAddress   *string  `form:"redeem_address" binding:"-"`
	ReferenceAmount *float64 `form:"reference_amount"`
}

// tx type: 3
type SendSto struct {
	TxVersion     string  `form:"transaction_version" binding:"required"`
	PubKey        *string `form:"pubkey" binding:"-"`
	Fee           float64 `form:"fee" binding:"gte=0.00001,required"`
	From          string  `form:"transaction_from" binding:"required"`
	PID           int64   `form:"propertyid" binding:"gte=1,required"`
	Amount        float64 `form:"amount" binding:"gt=0,required"`
	RedeemAddress *string `form:"redeem_address" binding:"-"`
	DisPID        *int64  `form:"distributionproperty"`
}

// tx type: 4
type SendAll struct {
	TxVersion       string   `form:"transaction_version" binding:"required"`
	PubKey          *string  `form:"pubkey" binding:"-"`
	Fee             float64  `form:"fee" binding:"gte=0.00001,required"`
	From            string   `form:"transaction_from" binding:"required"`
	To              string   `form:"transaction_to" binding:"required"`
	Eco             int64    `form:"ecosystem" binding:"eq=1,required"`
	RedeemAddress   *string  `form:"redeem_address"`
	ReferenceAmount *float64 `form:"reference_amount"`
}

// tx type: 50
type FixedIssuance struct {
	TxVersion   string  `form:"transaction_version" binding:"required"`
	PubKey      *string `form:"pubkey" binding:"-"`
	Fee         float64 `form:"fee" binding:"gte=0.00001,required"`
	From        string  `form:"transaction_from" binding:"required"`
	Eco         int64   `form:"ecosystem" binding:"eq=1,required"`
	Precision   int64   `form:"precision" binding:"max=8"`
	PrevPID     int64   `form:"previous_property_id" binding:"eq=0"`
	Category    string  `form:"property_category" binding:"required"`
	SubCategory string  `form:"property_subcategory" binding:"required"`
	Name        string  `form:"property_name" binding:"required"`
	Url         string  `form:"property_url" binding:"url,required"`
	Data        string  `form:"property_data" binding:"required"`
	TotalNumber string  `form:"number_properties" binding:"gt=0,required"`
}

// tx type: 51
type CrowdSaleIssuance struct {
	TxVersion   string  `form:"transaction_version" binding:"required"`
	PubKey      *string `form:"pubkey" binding:"-"`
	Fee         float64 `form:"fee" binding:"gte=0.00001,required"`
	From        string  `form:"transaction_from" binding:"required"`
	Eco         int64   `form:"ecosystem" binding:"eq=1,required"`
	Precision   int64   `form:"precision" binding:"max=8"`
	PrevPID     int64   `form:"previous_property_id" binding:"eq=0"`
	Category    string  `form:"property_category" binding:"required"`
	SubCategory string  `form:"property_subcategory" binding:"required"`
	Name        string  `form:"property_name" binding:"required"`
	Url         string  `form:"property_url" binding:"url,required"`
	Data        string  `form:"property_data" binding:"required"`
	DesiredPID  int64   `form:"currency_identifier_desired" binding:"eq=1,required"`
	Exchange    string  `form:"number_properties" binding:"gt=0,required"`
	DeadLine    int64   `form:"deadline" binding:"required"`
	EarlyBird   int64   `form:"earlybird_bonus" binding:"min=0,max=255"`
	TotalNumber string  `form:"total_number" binding:"gt=0,required"`
}

// tx type: 53
type CloseCrowdSale struct {
	TxVersion string  `form:"transaction_version" binding:"required"`
	PubKey    *string `form:"pubkey" binding:"-"`
	Fee       float64 `form:"fee" binding:"gte=0.00001,required"`
	From      string  `form:"transaction_from" binding:"required"`
	CID       int64   `form:"currency_identifier" binding:"gte=1,required"`
}

// tx type: 54
type ManagedIssuance struct {
	TxVersion   string  `form:"transaction_version" binding:"required"`
	PubKey      *string `form:"pubkey" binding:"-"`
	Fee         float64 `form:"fee" binding:"gte=0.00001,required"`
	From        string  `form:"transaction_from" binding:"required"`
	Eco         int64   `form:"ecosystem" binding:"eq=1,required"`
	Precision   int64   `form:"precision" binding:"max=8"`
	PrevPID     int64   `form:"previous_property_id"`
	Category    string  `form:"property_category" binding:"required"`
	SubCategory string  `form:"property_subcategory" binding:"required"`
	Name        string  `form:"property_name" binding:"required"`
	Url         string  `form:"property_url" binding:"url,required"`
	Data        string  `form:"property_data" binding:"required"`
}

// tx type: 55
type SendGrant struct {
	TxVersion string  `form:"transaction_version" binding:"required"`
	PubKey    *string `form:"pubkey" binding:"-"`
	Fee       float64 `form:"fee" binding:"gte=0.00001,required"`
	From      string  `form:"transaction_from" binding:"required"`
	CID       int64   `form:"currency_identifier" binding:"gte=1,required"`
	Amount    float64 `form:"amount" binding:"gt=0,required"`
	Note      *string `form:"note" binding:"-"`
}

// tx type: 56
type SendRevoke struct {
	TxVersion string  `form:"transaction_version" binding:"required"`
	PubKey    *string `form:"pubkey" binding:"-"`
	Fee       float64 `form:"fee" binding:"gte=0.00001,required"`
	From      string  `form:"transaction_from" binding:"required"`
	CID       int64   `form:"currency_identifier" binding:"gte=1,required"`
	Amount    float64 `form:"amount" binding:"gt=0,required"`
	Note      *string `form:"note" binding:"-"`
}

// tx type: 68
type BurnBCH struct {
	TxVersion     string  `form:"transaction_version" binding:"required"`
	PubKey        *string `form:"pubkey" binding:"-"`
	Fee           float64 `form:"fee" binding:"gte=0.00001,required"`
	From          string  `form:"transaction_from" binding:"required"`
	Amount        float64 `form:"amount_for_burn" binding:"min=1,max=21000000,required"`
	RedeemAddress *string `form:"redeem_address" binding:"-"`
}

// tx type: 70
type ChangeIssuer struct {
	TxVersion string  `form:"transaction_version" binding:"required"`
	PubKey    *string `form:"pubkey" binding:"-"`
	Fee       float64 `form:"fee" binding:"gte=0.00001,required"`
	From      string  `form:"transaction_from" binding:"required"`
	To        string  `form:"transaction_to" binding:"required"`
	CID       int64   `form:"currency_identifier" binding:"gte=1,required"`
}

type UserTx struct {
	FromAddress   string    `form:"transaction_from"`
	ToAddress     []string  `form:"address[]" binding:"required"`
	ToAmount      []float64 `form:"amount[]" binding:"required"`
	Fee           float64   `form:"fee" binding:"gte=0.00001,required"`
	RedeemAddress string    `form:"redeem_address"`
}

type FeeRate struct {
	Fast   float64
	Normal float64
	Slow   float64
}

type FreezeToken struct {
	From          string `form:"transaction_from" binding:"required"`
	PropertyId    int64  `form:"property_id" binding:"gte=1,required"`
	Amount        string `form:"amount"`
	FrozenAddress string `form:"frozen_address" binding:"-"`
}

type AddressBalance struct {
	Address          string
	BalanceAvailable *decimal.Decimal `json:"balance_available"`
	Status           int
}
