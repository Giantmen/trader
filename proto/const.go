package proto

const (
	HuobiN   = "HuobiN"
	HuobiO   = "HuobiO"
	Chbtc    = "Chbtc"
	Yunbi    = "Yunbi"
	Btctrade = "Btctrade"
	Btc38    = "Btc38"
	Jubi     = "Jubi"
	Bter     = "Bter"

	LocalTime = "2006-01-02 15:04:05"
	UTCTime   = "2006-01-02T15:04:05"
)

//手续费
const (
	FEE_Huobi_btc = 0.002
	FEE_Huobi_ltc
	FEE_Chbtc_btc
	FEE_Chbtc_ltc

	FEE_Chbtc_etc = 0.0005
	FEE_Huobi_etc //7月13日12:00-7月16日12:00 0.01%
	FEE_Yunbi_btc
	FEE_Btctrade_eth

	FEE_Yunbi_etc = 0.001
	FEE_Btctrade_etc

	FEE_Btctrade_btc = 0.0018
	FEE_Btctrade_ltc
)

const (
	BTC_CNY = "btc_cny"
	LTC_CNY = "ltc_cny"
	ETH_CNY = "eth_cny"
	ETH_BTC = "eth_btc"

	ETC_CNY = "etc_cny"
	ETC_BTC = "etc_btc"

	BTS_CNY = "bts_cny"
	BTS_BTC = "bts_btc"
	EOS_CNY = "eos_cny"
	SNT_CNY = "snt_cny"
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
