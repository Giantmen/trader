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
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
	//Timestamp int         `json:"timestamp"`
}

type MyAccount struct {
	Limit  float64 `json:"limit"`
	Result struct {
		Balance struct {
			BTC struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"BTC"`
			BTS struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"BTS"`
			CNY struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"CNY"`
			DAO struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"DAO"`
			ETC struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"ETC"`
			ETH struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"ETH"`
			LTC struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"LTC"`
		} `json:"balance"`
		Base struct {
			AuthGoogleEnabled    bool   `json:"auth_google_enabled"`
			AuthMobileEnabled    bool   `json:"auth_mobile_enabled"`
			TradePasswordEnabled bool   `json:"trade_password_enabled"`
			Username             string `json:"username"`
		} `json:"base"`
		Frozen struct {
			BTC struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"BTC"`
			BTS struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"BTS"`
			CNY struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"CNY"`
			DAO struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"DAO"`
			ETC struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"ETC"`
			ETH struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"ETH"`
			LTC struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Symbol   string  `json:"symbol"`
			} `json:"LTC"`
		} `json:"frozen"`
		NetAssets float64 `json:"netAssets"`
		P2p       struct {
			InBTC  float64 `json:"inBTC"`
			InCNY  float64 `json:"inCNY"`
			InDAO  float64 `json:"inDAO"`
			InETC  float64 `json:"inETC"`
			InETH  float64 `json:"inETH"`
			InLTC  float64 `json:"inLTC"`
			OutBTC float64 `json:"outBTC"`
			OutCNY float64 `json:"outCNY"`
			OutDAO float64 `json:"outDAO"`
			OutETC float64 `json:"outETC"`
			OutETH float64 `json:"outETH"`
			OutLTC float64 `json:"outLTC"`
		} `json:"p2p"`
		TotalAssets float64 `json:"totalAssets"`
	} `json:"result"`
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

type PlaceOrder struct {
	Code    int    `json:"code"`
	ID      string `json:"id"`
	Message string `json:"message"`
}
