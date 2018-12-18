package ws

type CoinType string

const (
	CoinBch      CoinType = "BCH"
	CoinWormhole CoinType = "Wormhole"
)

type Message struct {
	Address string      `json:"address"`
	Balance interface{} `json:"balance"`
	Symbol  CoinType    `json:"symbol"`
}

type BCHBalance struct {
	Confirmed   float64 `json:"confirmed"`
	Unconfirmed float64 `json:"unconfirmed"`
}
