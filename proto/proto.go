package proto

type DepthQuery struct {
	Size         int    `json:"size"`
	Depth        int    `json:"depth"`
	Bourse       string `validate:"required" json:"bourse"`
	CurrencyPair string `validate:"required" json:"currencypair"`
}

type Price struct {
	Buy  float64 `json:"buy,string"`
	Sell float64 `json:"sell,string"`
}

type Account struct {
	Bourse      string
	Asset       float64
	SubAccounts map[string]SubAccount
}

type SubAccount struct {
	Currency  string
	Available float64
	Forzen    float64
}

//cancel order
type CancelQuery struct {
	Bourse       string `validate:"required" json:"bourse"`
	OrderID      string `validate:"required" json:"orderid"`
	CurrencyPair string `validate:"required" json:"currencypair"`
}

//order
type OrderQuery struct {
	Bourse       string `validate:"required" json:"bourse"`
	Side         string `validate:"required" json:"side"`
	Amount       string `validate:"required" json:"amount"`
	Price        string `validate:"required" json:"price"`
	CurrencyPair string `validate:"required" json:"currencypair"`
}

type Order struct {
	OrderID      string  `validate:"required" json:"orderid"`
	OrderTime    string  `validate:"required" json:"ordertime"`
	Price        float64 `validate:"required" json:"price"`
	Amount       float64 `validate:"required" json:"amount"`
	DealedAmount float64 `validate:"required" json:"dealedamount"`
	Fee          float64 `validate:"required" json:"fee"`
	Status       string  `validate:"required" json:"status"`
	Currency     string  `validate:"required" json:"currency"`
	Side         string  `validate:"required" json:"side"`
}

//amount
type AmountQuery struct {
	Bourse   string   `validate:"required" json:"bourse"`
	Accounts []string `validate:"required" json:"accounts"`
}

type AmountReply struct {
	Bourse   string                `validate:"required" json:"bourse"`
	Asset    float64               `validate:"required" json:"asset"`
	Accounts map[string]SubAccount `json:"accounts"`
}
