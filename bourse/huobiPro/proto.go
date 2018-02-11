package huobiPro

import (
	"encoding/json"
)

type response struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Errmsg  string          `json:"err-msg"`
	Errcode string          `json:"err-code"`
}

type MyOrder struct {

}

type Depth struct {
	Tick Tick
	Status string
}
type Tick struct {
	Asks  [][]float64
	Bids  [][]float64
}