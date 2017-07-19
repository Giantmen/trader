package huobiN

import (
	"encoding/json"
	"fmt"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

var (
	API_URL = "https://be.huobi.com"
	//TICKER_URL     = "tickers/%s.json"
	DEPTH_URL = "http://be.huobi.com/market/kline?symbol=%s&period=1min"
	//USER_INFO_URL    = "members/me.json"
	CREATE_ORDER_API = "/v1/order/orders"
	PLACE_ORDER_API  = "/v1/order/orders/%d/place"
	//DELETE_ORDER_API = "order/delete.json"

)

type Huobi struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewHuobi(accessKey, secretKey string, timeout int) (*Huobi, error) {
	return &Huobi{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (h *Huobi) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf(DEPTH_URL, h.convertCurrencyPair(currencyPair), size)
	rep, err := util.Request("GET", url, "application/json", nil, nil, h.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Yunbi, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Huobi, currencyPair, err)
	}
	if body.Status != "ok" {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Huobi, currencyPair, body.Status)
	}

	var sellsum float64
	var sellprice float64
	var buysum float64
	var buyprice float64
	var len int = len(body.Tick.Asks)

	for i := 0; i < len; i++ {
		sellsum += body.Tick.Asks[i][1]
		if sellsum > float64(depth) {
			sellprice = body.Tick.Asks[i][0]
			break
		}
	}
	for i := 0; i < len; i++ {
		buysum += body.Tick.Bids[i][1]
		if buysum > float64(depth) {
			buyprice = body.Tick.Bids[i][0]
			break
		}

	}
	if sellsum > float64(depth) && buysum > float64(depth) {
		return &proto.Price{
			Sell: sellprice,
			Buy:  buyprice,
		}, nil
	}
	return nil, fmt.Errorf("sum not enough %v %v", sellsum, depth)
}

func (h *Huobi) convertCurrencyPair(currencyPair string) string {
	switch currencyPair {
	case proto.ETH_CNY:
		return "ethcny"
	case proto.ETC_CNY:
		return "etccny"
	default:
		return ""
	}
}
