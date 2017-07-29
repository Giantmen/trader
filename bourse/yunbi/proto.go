package yunbi

type Market struct {
	At     uint64 `json:"at"`
	Ticker Ticker `json:"ticker"`
}

type Ticker struct {
	Buy  float64 `json:"buy,string"`
	Sell float64 `json:"sell,string"`
	Low  float64 `json:"low,string"`
	High float64 `json:"high,string"`
	Last float64 `json:"last,string"`
	Vol  float64 `json:"vol,string"`
}

//depth
type Depth struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	//Timestamp int         `json:"timestamp"`
}

type MyAccount struct {
	Accounts []struct {
		Balance  string `json:"balance"`
		Currency string `json:"currency"`
		Locked   string `json:"locked"`
	} `json:"accounts"`
	Activated bool   `json:"activated"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Sn        string `json:"sn"`
}

type MyOrder struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	AvgPrice        string `json:"avg_price"`
	CreatedAt       string `json:"created_at"`
	ExecutedVolume  string `json:"executed_volume"`
	ID              int    `json:"id"`
	Market          string `json:"market"`
	Price           string `json:"price"`
	RemainingVolume string `json:"remaining_volume"`
	Side            string `json:"side"`
	State           string `json:"state"`
	Trades          []struct {
		CreatedAt string `json:"created_at"`
		ID        int    `json:"id"`
		Market    string `json:"market"`
		Price     string `json:"price"`
		Side      string `json:"side"`
		Volume    string `json:"volume"`
	} `json:"trades"`
	TradesCount int    `json:"trades_count"`
	Volume      string `json:"volume"`
}
