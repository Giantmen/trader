package poloniex

type Depth struct {
	Asks  interface{}
	Bids  interface{}
	Error string
}

// 下单的返回值

type PlaceOrder struct {
	OrderNumber string `json:"orderNumber"`
	Error       string `json:"error"`
	// ResultingTrades []struct {
	// 	Amount  string `json:"amount"`
	// 	Date    string `json:"date"`
	// 	Rate    string `json:"rate"`
	// 	Total   string `json:"total"`
	// 	TradeID string `json:"tradeID"`
	// 	Type    string `json:"type"`
	// } `json:"resultingTrades"`
}

// 获取账户的返回值
type MyAccount struct {
	Exchange struct {
		BTC string `json:"BTC"`
		LTC string `json:"LTC"`
		ETC string `json:"ETC"`
		ETH string `json:"ETH"`
	} `json:"exchange"`
}

type MyOrder struct {
	Orders []Order `json:"myOrders"`
	Error  string  `json:"error"`
}

type Order struct {
	Amount        string `json:"amount"`
	CurrencyPair  string `json:"currencyPair"`
	Date          string `json:"date"`
	Fee           string `json:"fee"`
	GlobalTradeID int64  `json:"globalTradeID"`
	Rate          string `json:"rate"`
	Total         string `json:"total"`
	TradeID       int64  `json:"tradeID"`
	Type          string `json:"type"`
}

type OpenOrder struct {
	OrderNumber string
	Type        string
	Rate        string
	Amount      string
	Total       string
}
