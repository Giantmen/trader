package bitfinex

type Depth struct {
	Asks  []SubDepth
	Bids  []SubDepth
	Error string
}

type SubDepth struct {
	Price string
	Amount string
	Timestamp string
}

type MyOrder struct {
	Price      float64
	Amount     float64
	AvgPrice   float64
	DealAmount float64
	Fee        float64
	OrderID2   string
	OrderID    int
	OrderTime  int
	Status     int
	Currency   string
	Side       int
}

type MyAccount struct {
	Type string
	Currency string
	Amount string
	Available string
}
