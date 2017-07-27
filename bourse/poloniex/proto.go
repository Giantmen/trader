package poloniex

type Depth struct {
	Asks [][]interface{} `json:"asks"`
	Bids [][]interface{} `json:"bids"`
	//IsFrozen string          `json:"isFrozen"`
	//Seq      int             `json:"seq"`
}

type MyAccount struct {
}

type MyOrder struct {
}
