package poloniex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

var (
	API_URL          = "https://poloniex.com/public?"
	TICKER_URL       = "tickers/%s.json"
	DEPTH_URL        = "command=returnOrderBook&currencyPair=%s&depth=%d"
	USER_INFO_URL    = "members/me.json"
	GET_ORDER_API    = "order.json"
	DELETE_ORDER_API = "order/delete.json"
	PLACE_ORDER_API  = "orders.json"
)

type Poloniex struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewPoloniex(accessKey, secretKey string, timeout int) (*Poloniex, error) {
	return &Poloniex{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (poloniex *Poloniex) GetTicker(currencyPair string) (float64, error) {
	return 0, nil
}

func (poloniex *Poloniex) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := API_URL + fmt.Sprintf(DEPTH_URL, strings.ToUpper(currencyPair), size)
	rep, err := util.Request("GET", url, "application/json", nil, nil, 4)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Poloniex, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Poloniex, currencyPair, err)
	}

	var sellsum float64
	var buysum float64
	var sellprice float64
	var buyprice float64
	var len int = len(body.Asks)
	for i := 0; i < len; i++ {
		price, err := strconv.ParseFloat((body.Asks[i][0]).(string), 64)
		if err != nil {
			continue
		}
		sum := price * ((body.Asks[i][1]).(float64))
		sellsum += sum
		if sellsum > float64(depth) {
			sellprice = price
			break
		}
	}

	for i := 0; i < len; i++ {
		price, err := strconv.ParseFloat((body.Bids[i][0]).(string), 64)
		if err != nil {
			continue
		}
		sum := price * ((body.Bids[i][1]).(float64))
		buysum += sum
		if buysum > float64(depth) {
			buyprice = price
			break
		}
	}

	if sellsum > float64(depth) && buysum > float64(depth) {
		price := proto.Price{
			Sell: sellprice,
			Buy:  buyprice,
		}
		return &price, nil
	}
	return nil, fmt.Errorf("sum not enough sell:%v buy:%v depth:%v", sellsum, buysum, depth)
}

func (poloniex *Poloniex) GetAccount() (*proto.Account, error) {
	return nil, nil
}

func (poloniex *Poloniex) placeOrder(side int, amount, price, currencyPair string) (*proto.Order, error) {
	return nil, nil
}

func (poloniex *Poloniex) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return poloniex.placeOrder(proto.BUY_N, amount, price, currencyPair)
}

func (poloniex *Poloniex) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return poloniex.placeOrder(proto.SELL_N, amount, price, currencyPair)
}

func (poloniex *Poloniex) CancelOrder(orderId, currencyPair string) (bool, error) {
	return false, nil
}

func (poloniex *Poloniex) parseOrder(myorder *MyOrder) (*proto.Order, error) {
	return nil, nil
}

func (poloniex *Poloniex) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	return nil, nil
}

func (poloniex *Poloniex) buildPostForm(postForm *url.Values) error {
	return nil
}
