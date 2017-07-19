package jubi

type Depth struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
	//Timestamp int         `json:"timestamp"`
}

type MyOrder struct {
	Currency    string  `json:"currency"`
	Fees        float64 `json:"fees"`
	ID          string  `json:"id"`
	Price       float64 `json:"price"`
	Status      float64 `json:"status"`
	TotalAmount float64 `json:"total_amount"`
	TradeAmount float64 `json:"trade_amount"`
	TradeDate   float64 `json:"trade_date"`
	TradeMoney  float64 `json:"trade_money"`
	TradePrice  float64 `json:"trade_price"`
	Type        int     `json:"type"`
}
