package poloniex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

const (
	PUBLICAPI = "https://poloniex.com/public"
	TRADEAPI  = "https://poloniex.com/tradingApi"
)

const (
	BUY              = "buy"
	SELL             = "sell"
	ORDERBOOK        = "returnOrderBook"
	ORDERTRADES      = "returnOrderTrades"
	OPENORDERS       = "returnOpenOrders"
	CANCLEORDER      = "cancelOrder"
	COMPLETEBALANCES = "returnAvailableAccountBalances"
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

func (p *Poloniex) GetTicker(currencyPair string) (float64, error) {
	return 0.0, nil
}

// 获取满足某个深度的价格
func (p *Poloniex) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf("%s?command=%s&currencyPair=%s&depth=%d", PUBLICAPI, ORDERBOOK, strings.ToUpper(currencyPair), size)
	rep, err := util.Request("GET", url, "", nil, nil, p.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Poloniex, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Poloniex, currencyPair, err)
	}
	if body.Error != "" {
		return nil, fmt.Errorf("err: %s", body.Error)
	}

	asks, _ := body.Asks.([]interface{})
	bids, _ := body.Bids.([]interface{})
	buyPrice, buySum, err := priceOfDepth(asks, depth)
	if err != nil {
		return nil, err
	}
	sellPrice, sellSum, err := priceOfDepth(bids, depth)
	if err != nil {
		return nil, err
	}
	return &proto.Price{
		Buy:     buyPrice,
		Sell:    sellPrice,
		Sellnum: sellSum,
		Buynum:  buySum,
	}, nil
}

func priceOfDepth(terms []interface{}, depth float64) (float64, float64, error) {
	var sum float64 = 0.0
	for _, term := range terms {
		entry, _ := term.([]interface{})
		p := entry[0].(string)
		c := entry[1].(float64)
		price, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0.0, 0.0, err
		}
		total := price * c
		sum += total
		if sum > float64(depth) {
			return price, sum, nil
		}
	}
	return 0.0, 0.0, fmt.Errorf("has no enough depth sum:%v depth:%v", sum, depth)
}

func (p *Poloniex) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	v := url.Values{}
	v.Set("command", side)
	v.Set("currencyPair", strings.ToUpper(currencyPair))
	v.Set("rate", price)
	v.Set("amount", amount)

	sign, err := p.buildPostForm(&v)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Key", p.accessKey)
	header.Set("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", TRADEAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, fmt.Errorf("request %v err:%v", side, err)
	}
	tresp := new(PlaceOrder)
	err = json.Unmarshal(resp, &tresp)
	if err != nil {
		return nil, err
	}
	if tresp.Error != "" {
		return nil, fmt.Errorf("%s err:%s", side, tresp.Error)
	}
	return &proto.Order{
		OrderID:      tresp.OrderNumber,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Price:        float64(0),
		Amount:       float64(0),
		DealedAmount: float64(0),
		Fee:          float64(0),
		Status:       proto.ORDER_UNFINISH,
		Currency:     currencyPair,
		Side:         side,
	}, nil
}

func (p *Poloniex) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return p.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (p *Poloniex) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return p.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (p *Poloniex) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	v := url.Values{}
	v.Set("command", ORDERTRADES)
	v.Set("orderNumber", orderId)
	sign, _ := p.buildPostForm(&v)
	header := http.Header{}
	header.Set("Key", p.accessKey)
	header.Set("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", TRADEAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}

	myOrder := MyOrder{}
	err = json.Unmarshal(resp, &myOrder)
	if err != nil {
		return nil, err
	}
	if myOrder.Error != "" {
		return p.getUnfinishOrders(orderId, currencyPair)
	}
	if len(myOrder.Orders) == 0 {
		return nil, fmt.Errorf("request GetOneOrder err %s", string(resp))
	}

	order := new(proto.Order)
	order.OrderID = orderId
	order.Currency = currencyPair
	if len(myOrder.Orders) > 0 {
		order.OrderTime = myOrder.Orders[0].Date
	}
	var amounts float64
	for _, myorder := range myOrder.Orders {
		amount, _ := strconv.ParseFloat(myorder.Total, 64)
		amounts += amount
	}
	order.DealedAmount = amounts
	return order, nil
}

func (p *Poloniex) getUnfinishOrders(orderId, currencyPair string) (*proto.Order, error) {
	v := url.Values{}
	v.Set("command", OPENORDERS)
	v.Set("currencyPair", strings.ToUpper(currencyPair))
	sign, err := p.buildPostForm(&v)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Key", p.accessKey)
	header.Set("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", TRADEAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}

	openOrder := make([]OpenOrder, 1)
	err = json.Unmarshal(resp, &openOrder)
	if err != nil {
		return nil, err
	}
	for _, order := range openOrder {
		if order.OrderNumber == orderId {
			price, _ := strconv.ParseFloat(order.Rate, 64)
			amount, _ := strconv.ParseFloat(order.Amount, 64)
			return &proto.Order{
				OrderID:   orderId,
				Side:      order.Type,
				Price:     price,
				Amount:    amount,
				Currency:  currencyPair,
				OrderTime: time.Now().Format(proto.LocalTime),
				Status:    proto.ORDER_UNFINISH,
			}, nil
		}
	}
	return nil, fmt.Errorf("can not find orderid:%s", orderId)
}

func (p *Poloniex) CancelOrder(orderId, currencypair string) (bool, error) {
	v := url.Values{}
	v.Set("command", CANCLEORDER)
	v.Set("orderNumber", orderId)
	sign, err := p.buildPostForm(&v)
	if err != nil {
		return false, err
	}

	header := http.Header{}
	header.Set("Key", p.accessKey)
	header.Set("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", TRADEAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return false, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(resp, &result)
	if err != nil || result["error"] != nil {
		return false, err
	}

	success := int(result["success"].(float64))
	if success != 1 {
		return false, nil
	}
	return true, nil
}

func (p *Poloniex) GetAccount() (*proto.Account, error) {
	v := url.Values{}
	v.Set("command", COMPLETEBALANCES)
	sign, err := p.buildPostForm(&v)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Set("Key", p.accessKey)
	header.Set("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", TRADEAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}
	myaccount := MyAccount{}
	if resp == nil || string(resp) == "[]" {
		goto ACCOUNT
	}
	err = json.Unmarshal(resp, &myaccount)
	if err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v", err)
	}

ACCOUNT:
	account := proto.Account{}
	account.Asset = 0
	account.Bourse = strings.ToLower(proto.Poloniex)
	account.SubAccounts = make(map[string]proto.SubAccount)

	//btc
	subAcc := proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Exchange.BTC, 64)
	subAcc.Forzen = 0
	subAcc.Currency = proto.BTC
	account.SubAccounts[subAcc.Currency] = subAcc
	//etc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Exchange.ETC, 64)
	subAcc.Forzen = 0
	subAcc.Currency = proto.ETC
	account.SubAccounts[subAcc.Currency] = subAcc
	//eth
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Exchange.ETH, 64)
	subAcc.Forzen = 0
	subAcc.Currency = proto.ETH
	account.SubAccounts[subAcc.Currency] = subAcc
	//ltc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Exchange.LTC, 64)
	subAcc.Forzen = 0
	subAcc.Currency = proto.LTC
	account.SubAccounts[subAcc.Currency] = subAcc

	return &account, nil
}

func (p *Poloniex) buildPostForm(v *url.Values) (string, error) {
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	payload := v.Encode()
	sign, err := util.SHA512Sign(p.secretKey, payload)
	if err != nil {
		return "", err
	}
	return sign, nil
}
