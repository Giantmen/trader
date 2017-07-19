package btctrade

//depth
type Depth struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	//Result string     `json:"result"`
}

type MyAccount struct {
	Asset       string `json:"asset"`
	BtcBalance  string `json:"btc_balance"`
	BtcReserved string `json:"btc_reserved"`
	CnyBalance  string `json:"cny_balance"`
	CnyReserved string `json:"cny_reserved"`
	// DogeBalance  string `json:"doge_balance"`
	// DogeReserved string `json:"doge_reserved"`
	EtcBalance  string `json:"etc_balance"`
	EtcReserved string `json:"etc_reserved"`
	EthBalance  string `json:"eth_balance"`
	EthReserved string `json:"eth_reserved"`
	LtcBalance  string `json:"ltc_balance"`
	LtcReserved string `json:"ltc_reserved"`
	Moflag      string `json:"moflag"`
	Nameauth    int    `json:"nameauth"`
	UID         int    `json:"uid"`
	// YbcBalance   string `json:"ybc_balance"`
	// YbcReserved  string `json:"ybc_reserved"`
}

type OrderReply struct {
	ID      string `json:"id"`
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type MyOrder struct {
	Result            bool    `json:"result"`
	Message           string  `json:"message"`
	AmountOriginal    float64 `json:"amount_original"`
	AmountOutstanding float64 `json:"amount_outstanding"`
	Datetime          string  `json:"datetime"`
	ID                int64   `json:"id"`
	Price             float64 `json:"price"`
	Status            string  `json:"status"`
	Trades            struct {
		AvgPrice  float64 `json:"avg_price"`
		SumMoney  float64 `json:"sum_money"`
		SumNumber float64 `json:"sum_number"`
	} `json:"trades"`
	Type string `json:"type"`
}
