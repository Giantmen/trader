package chbtc

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

//depth
type Depth struct {
	Asks  [][]float64 `json:"asks"`
	Bids  [][]float64 `json:"bids"`
	Error string      `json:"error"`
	//Timestamp int         `json:"timestamp"`
}

type MyAccount struct {
	Result struct {
		Balance struct {
			BCC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"BCC"`
			BTC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"BTC"`
			BTS struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"BTS"`
			CNY struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"CNY"`
			EOS struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"EOS"`
			ETC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"ETC"`
			ETH struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"ETH"`
			LTC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"LTC"`
			QTUM struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"QTUM"`
		} `json:"balance"`
		Base struct {
			AuthGoogleEnabled    bool   `json:"auth_google_enabled"`
			AuthMobileEnabled    bool   `json:"auth_mobile_enabled"`
			TradePasswordEnabled bool   `json:"trade_password_enabled"`
			Username             string `json:"username"`
		} `json:"base"`
		Frozen struct {
			BCC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"BCC"`
			BTC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"BTC"`
			BTS struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"BTS"`
			CNY struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"CNY"`
			EOS struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"EOS"`
			ETC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"ETC"`
			ETH struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"ETH"`
			LTC struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"LTC"`
			QTUM struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
				Symbol   string `json:"symbol"`
			} `json:"QTUM"`
		} `json:"frozen"`
		NetAssets   string `json:"netAssets"`
		TotalAssets string `json:"totalAssets"`
	} `json:"result"`
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
	Code    int    `json:"code"`
	ID      string `json:"id"`
	Message string `json:"message"`
}
