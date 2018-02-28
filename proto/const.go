package proto

import "strings"

const (
	Bittrex  = "Bittrex"
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
	case "bittrex":
		return 0.0025
	case "huobi_btc", "huobi_ltc", "chbtc_btc", "chbtc_ltc":
		return 0.002
	case "bter_snt", "bter_omg", "bter_pay", "bter_btm":
		return 0.001
	case "yunbi_etc", "yunbi_eth", "yunbi_snt", "yunbi_omg", "yunbi_pay",
		"btctrade_etc", "huobi_etc", "huobi_eth":
		return 0.001
	case "yunbi_btc", "btctrade_eth":
		return 0.0005
	case "chbtc_etc", "chbtc_eth":
		return 0.00046
	default:
		return 0
	}
}

const (
	BTC_CNY = "btc_cny"
	LTC_CNY = "ltc_cny"
	ETH_CNY = "eth_cny"
	ETC_CNY = "etc_cny"
<<<<<<< HEAD
=======

>>>>>>> 53dbc403f087b5c6d26a43c9cd9963b934690dfd
	BTS_CNY = "bts_cny"
	EOS_CNY = "eos_cny"
	SNT_CNY = "snt_cny"
	OMG_CNY = "omg_cny"
	PAY_CNY = "pay_cny"
	BTM_CNY = "btm_cny"

<<<<<<< HEAD
	BTC_LTC = "btc_ltc"
	BTC_ETH = "btc_eth"
	BTC_ETC = "btc_etc"
	BTC_BTS = "btc_bts"
	BTC_BTM = "btc_btm"
	BTC_EOS = "btc_eos"
	BTC_SNT = "btc_snt"
	BTC_OMG = "btc_omg"
	BTC_PAY = "btc_pay"
	BTC_CVC = "btc_cvc"
	BTC_SC  = "btc_sc"
=======
	ETH_BTC = "eth_btc"
	LTC_BTC = "ltc_btc"
	EOS_BTC = "eos_btc"
	NEO_BTC = "neo_btc"

>>>>>>> 53dbc403f087b5c6d26a43c9cd9963b934690dfd
)
const (
	CNY = "cny"
	BTC = "btc"
	LTC = "ltc"
	ETH = "eth"
	ETC = "etc"
	BTS = "bts"
	BTM = "btm"
	EOS = "eos"
	SNT = "snt"
	OMG = "omg"
	PAY = "pay"
<<<<<<< HEAD
	CVC = "cvc"
	SC  = "sc"
=======
	BTM = "btm"
	NEO = "neo"
>>>>>>> 53dbc403f087b5c6d26a43c9cd9963b934690dfd
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
