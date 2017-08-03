package huobiN

type Ticker struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Tick   struct {
		Amount float64 `json:"amount"`
		Close  float64 `json:"close"`
		Count  int     `json:"count"`
		High   float64 `json:"high"`
		ID     int     `json:"id"`
		Low    float64 `json:"low"`
		Open   float64 `json:"open"`
		Vol    float64 `json:"vol"`
	} `json:"tick"`
	Ts int `json:"ts"`
}

type Depth struct {
	// Ch     string `json:"ch"`
	//Ts int `json:"ts"`
	Status string `json:"status"`
	Tick   struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
		//Ts      int         `json:"ts"`
		//Version int         `json:"version"`
	} `json:"tick"`
}

type MyAccount struct {
	Data struct {
		//ID   int `json:"id"`
		List []struct {
			Balance  string `json:"balance"`
			Currency string `json:"currency"`
			Type     string `json:"type"`
		} `json:"list"`
		State string `json:"state"`
		Type  string `json:"type"`
	} `json:"data"`
	Status   string `json:"status"`
	Err_code string `json:"err-code"`
	Err_msg  string `json:"err-msg"`
}

type StatusOrder struct {
	Data     interface{} `json:"data"`
	Err_code string      `json:"err-code"`
	Err_msg  string      `json:"err-msg"`
	Status   string      `json:"status"`
}

type CreateOrder struct {
	AccountId string `json:"account-id"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Source    string `json:"source"`
	Symbol    string `json:"symbol"`
	Type      string `json:"type"`
}

type MyOrder struct {
	Data struct {
		Amount            string `json:"amount"`
		Created_at        int64  `json:"created-at"`
		Field_amount      string `json:"field-amount"`
		Field_cash_amount string `json:"field-cash-amount"`
		Field_fees        string `json:"field-fees"`
		ID                int64  `json:"id"`
		Price             string `json:"price"`
		State             string `json:"state"`
		Symbol            string `json:"symbol"`
		Type              string `json:"type"`
		// Account_id        int    `json:"account-id"`
		// Canceled_at       int64    `json:"canceled-at"`
		// Finished_at       int    `json:"finished-at"`
		// Source            string `json:"source"`
	} `json:"data"`
	Status   string `json:"status"`
	Err_code string `json:"err-code"`
	Err_msg  string `json:"err-msg"`
}
