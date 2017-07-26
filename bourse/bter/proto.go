package bter

type Market struct {
	Last   float64 `json:"last"`
	Result string  `json:"result"`
}

type Depth struct {
	Asks    [][]float64 `json:"asks"`
	Bids    [][]float64 `json:"bids"`
	Result  string      `json:"result"`
	Message string      `json:"message"`
}

type MyAccount struct {
	Available struct {
		BTC string `json:"BTC"`
		CNY string `json:"CNY"`
		LTC string `json:"LTC"`
		SNT string `json:"SNT"`
		ETC string `json:"ETC"`
		ETH string `json:"ETH"`
	} `json:"available"`
	Locked struct {
		BTC string `json:"BTC"`
		CNY string `json:"CNY"`
		LTC string `json:"LTC"`
		SNT string `json:"SNT"`
		ETC string `json:"ETC"`
		ETH string `json:"ETH"`
	} `json:"locked"`
	Result string `json:"result"`
}

type PlaceOrder struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	OrderNumber int64  `json:"orderNumber"`
	Result      string `json:"result"`
}

type MyOrder struct {
	Code    int    `json:"code"`
	Elapsed string `json:"elapsed"`
	Message string `json:"message"`
	Order   struct {
		CurrencyPair  string      `json:"currencyPair"`
		FilledAmount  float64     `json:"filledAmount"`
		FilledRate    interface{} `json:"filledRate"`
		InitialAmount string      `json:"initialAmount"`
		InitialRate   float64     `json:"initialRate"`
		OrderNumber   string      `json:"orderNumber"`
		Status        string      `json:"status"`
		Type          string      `json:"type"`
		//Amount       string `json:"amount"`
		//Fee          string `json:"fee"`
		//FeeCurrency  string `json:"feeCurrency"`
		//FeePercentage float64 `json:"feePercentage"` 0.18
		//FeeValue      string      `json:"feeValue"`
		//Rate          float64     `json:"rate"`
		//Timestamp     string  `json:"timestamp"`

	} `json:"order"`
	Result string `json:"result"`
}

type CancelOrder struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	//Result  interface{} `json:"result"`
}
