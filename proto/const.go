package proto

import "strings"

const (
	Poloniex = "Poloniex"
	Btctrade = "Btctrade"
	HuobiN   = "HuobiN"
	HuobiO   = "HuobiO"
	Chbtc    = "Chbtc"
	Yunbi    = "Yunbi"
	Btc38    = "Btc38"
	Jubi     = "Jubi"
	Bter     = "Bter"
	Binance = "Binance"
	Bitfinex = "Bitfinex"
	Huobipro = "Huobipro"

	LocalTime = "2006-01-02 15:04:05"
	UTCTime   = "2006-01-02T15:04:05"
)

func ConvertFee(brouse string) float64 {
	switch strings.ToLower(brouse) {
	case "huobi_btc", "huobi_ltc", "chbtc_btc", "chbtc_ltc":
		return 0.002
	case "yunbi_btc", "btctrade_eth":
		return 0.0005
	case "chbtc_etc", "chbtc_eth":
		return 0.00046
	case "bter_snt", "bter_omg", "bter_pay", "bter_btm":
		return 0.001
	case "yunbi_etc", "yunbi_eth", "yunbi_snt", "yunbi_omg", "yunbi_pay",
		"btctrade_etc", "huobi_etc", "huobi_eth":
		return 0.001
	default:
		return 0
	}
}

const (
	BTC_CNY = "btc_cny"
	LTC_CNY = "ltc_cny"
	ETH_CNY = "eth_cny"
	ETC_CNY = "etc_cny"

	BTS_CNY = "bts_cny"
	EOS_CNY = "eos_cny"

	SNT_CNY = "snt_cny"
	OMG_CNY = "omg_cny"
	PAY_CNY = "pay_cny"
	BTM_CNY = "btm_cny"

	ETH_BTC = "eth_btc"
	LTC_BTC = "ltc_btc"
	EOS_BTC = "eos_btc"
	NEO_BTC = "neo_btc"

)
const (
	CNY = "cny"
	BTC = "btc"
	LTC = "ltc"
	ETH = "eth"
	ETC = "etc"
	BTS = "bts"
	EOS = "eos"
	SNT = "snt"
	OMG = "omg"
	PAY = "pay"
	BTM = "btm"
	NEO = "neo"
)

const (
	BUY         = "buy"
	SELL        = "sell"
	BUY_N       = 1
	SELL_N      = 0
	BUY_MARKET  = "buy_market"
	SELL_MARKET = "sell_market"
)

const (
	ORDER_UNFINISH    = "UNFINISH"
	ORDER_PART_FINISH = "PART_FINISH"
	ORDER_FINISH      = "FINISH"
	ORDER_CANCEL      = "CANCEL"
	ORDER_REJECT      = "REJECT"
	ORDER_CANCEL_ING  = "CANCEL_ING"
)
