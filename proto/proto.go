package proto

type DepthQuery struct {
	Size     int     `json:"size"`
	Depth    float64 `json:"depth"`
	Bourse   string  `validate:"required" json:"bourse"`
	Currency string  `validate:"required" json:"currency"`
}

type Price struct {
	Buy     float64 `json:"buy,string"`
	Sell    float64 `json:"sell,string"`
	Buynum  float64 `json:"buynum,string"`
	Sellnum float64 `json:"sellnum,string"`
}

type Account struct {
	Bourse      string                `validate:"required" json:"bourse"`
	Asset       float64               `validate:"required" json:"asset"`
	SubAccounts map[string]SubAccount `json:"subAccounts"`
}

type SubAccount struct {
	Currency  string
	Available float64
	Forzen    float64
}

//cancel order
type OneOrderQuery struct {
	Bourse   string `validate:"required" json:"bourse"`
	OrderID  string `validate:"required" json:"orderid"`
	Currency string `validate:"required" json:"currency"`
}

//order
type OrderQuery struct {
	Bourse   string `validate:"required" json:"bourse"`
	Amount   string `validate:"required" json:"amount"`
	Price    string `validate:"required" json:"price"`
	Currency string `validate:"required" json:"currency"`
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

//account
type AccountQuery struct {
	Bourse string `validate:"required" json:"bourse"`
	//Accounts []string `validate:"required" json:"accounts"`
}
