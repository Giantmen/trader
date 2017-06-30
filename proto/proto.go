package proto

type Order struct {
	OrderID      string
	OrderTime    string
	Price        float64
	Amount       float64
	DealedAmount float64
	Fee          float64
	Status       string
	Currency     string
	Side         string
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
