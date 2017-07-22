package huobiO

type Ticker struct {
	Ticker struct {
		Buy    float64 `json:"buy"`
		High   float64 `json:"high"`
		Last   float64 `json:"last"`
		Low    float64 `json:"low"`
		Open   float64 `json:"open"`
		Sell   float64 `json:"sell"`
		Symbol string  `json:"symbol"`
		Vol    float64 `json:"vol"`
	} `json:"ticker"`
	Time string `json:"time"`
}
