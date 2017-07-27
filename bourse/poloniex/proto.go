package poloniex

type OrderBook struct {
	Asks interface{}
	Bids interface{}
}

type OpenOrder struct {
	OrderNumber string
	Type        string
	Rate        string
	Amount      string
	Total       string
}

// 下单的返回值
type TradeResponse struct {
	OrderNumber     string
	ResultingTrades []ResultingTrade
}

type ResultingTrade struct {
	Amount  string
	Date    string
	Rate    string
	Total   string
	TradeID string
	Type    string
}

// 获取账户的返回值
type AccountResponse map[string]SubAccount

type SubAccount struct {
	Available string
	OnOrders  string
	BtcValue  string
}

type OrderResponse []TradeTerm

type TradeTerm struct {
	GlobalTradeID int64
	TradeID       int64
	CurrencyPair  string
	Type          string
	Rate          string
	Amount        string
	Total         string
	Fee           string
	Date          string
}

// 操作
type Command string

func (c Command) String() string {
	return string(c)
}
