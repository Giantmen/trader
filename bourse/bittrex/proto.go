package bittrex

import "encoding/json"

type Response struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type Depth struct {
	Buy []struct {
		Quantity float64 `json:"Quantity"`
		Rate     float64 `json:"Rate"`
	} `json:"buy"`
	Sell []struct {
		Quantity float64 `json:"Quantity"`
		Rate     float64 `json:"Rate"`
	} `json:"sell"`
}

type Uuid struct {
	Id string `json:"uuid"`
}
