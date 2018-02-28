package binance

type Respons struct {
	Code    int
	Message string
}

type Market struct {
	Date   string `json:"date"`
	Ticker Ticker `json:"ticker"`
}
type Ticker struct {
	Buy  string `json:"buy"`
	High string `json:"high"`
	Last string `json:"last"`
	Low  string `json:"low"`
	Sell string `json:"sell"`
	Vol  string `json:"vol"`
}

type Depth struct {
	Asks  interface{}
	Bids  interface{}
	Error string
}

type SubAccount struct {
	Currency string
	Amount,
	ForzenAmount,
	LoanAmount float64
}

type Account struct {
	Exchange    string
	Asset       float64 //总资产
	NetAsset    float64 //净资产
	SubAccounts map[string]SubAccount
}

type MyOrder struct {
	Code        int     `json:"code"`
	Message     string  `json:"message"`
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

type PlaceOrder struct {
	Code    int
	Msg string
	Symbol string
	OrderId int64
	ClientOrderId string
	OrderTime int64
	Price string
	OrigQty string
	ExecutedQty string
	Status string
	TimeInForce string
	Type string
	Side string
}
