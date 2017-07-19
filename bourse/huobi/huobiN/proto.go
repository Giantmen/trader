<<<<<<< HEAD
package huobin

type Test struct{}
=======
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
		Asks    [][]float64 `json:"asks"`
		Bids    [][]float64 `json:"bids"`
		Ts      int         `json:"ts"`
		Version int         `json:"version"`
	} `json:"tick"`
}

//Response
type OrderResponse struct {
	Data   int    `json:"data"`
	Status string `json:"status"`
}

type CreateOrder struct {
	Account_id string `json:"account-id"`
	Amount     string `json:"amount"`
	Price      string `json:"price"`
	Source     string `json:"source"`
	Symbol     string `json:"symbol"`
	Type       string `json:"type"`
}

type Order struct {
	AccessKeyID      string `json:"AccessKeyId"`
	SignatureMethod  string `json:"SignatureMethod"`
	SignatureVersion int    `json:"SignatureVersion"`
	Timestamp        string `json:"Timestamp"`
	Order_id         string `json:"order-id"`
}
>>>>>>> 061303c1cd50f8136b235d820cd78db03cff40df
