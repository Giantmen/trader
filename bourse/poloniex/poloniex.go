package poloniex

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

const (
	baseURL   = "https://poloniex.com/"
	tradeAPI  = "https://poloniex.com/tradingApi"
	publicAPI = "https://poloniex.com/public"
)

const (
	BUY              = "buy"
	SELL             = "sell"
	ORDERBOOK        = "returnOrderBook"
	ORDERTRADES      = "returnOrderTrades"
	OPENORDERS       = "returnOpenOrders"
	CANCLEORDER      = "cancelOrder"
	COMPLETEBALANCES = "returnCompleteBalances"
)

type Poloniex struct {
	accessKey string
	secretkey string
	timeout   int
}

func NewPoloniex(accessKey, secretKey string, timeout int) (*Poloniex, error) {
	return &Poloniex{
		accessKey: accessKey,
		secretkey: secretKey,
		timeout:   timeout,
	}, nil
}

func (p *Poloniex) GetTicker(currencyPair string) (float64, error) {
	return 0.0, nil
}

// 获取满足某个深度的价格
func (p *Poloniex) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf("%s?command=%s&currencyPair=%s&depth=%d", publicAPI, ORDERBOOK, strings.ToUpper(currencyPair), size)
	log.Println(url)
	resp, err := util.Request("GET", url, "", nil, nil, p.timeout)
	if err != nil {
		return nil, err
	}

	ob := new(OrderBook)
	err = json.Unmarshal(resp, ob)
	if err != nil {
		return nil, err
	}

	asks, _ := ob.Asks.([]interface{})
	bids, _ := ob.Bids.([]interface{})
	buyPrice, err := priceOfDepth(asks, depth)
	if err != nil {
		return nil, err
	}
	sellPrice, err := priceOfDepth(bids, depth)
	if err != nil {
		return nil, err
	}
	return &proto.Price{
		Buy:  buyPrice,
		Sell: sellPrice,
	}, nil
}

func priceOfDepth(terms []interface{}, depth float64) (float64, error) {
	var d float64 = 0.0
	for _, term := range terms {
		entry, _ := term.([]interface{})
		p := entry[0].(string)
		c := entry[1].(float64)
		pf, err := strconv.ParseFloat(p, 32)
		if err != nil {
			return 0.0, err
		}
		total := pf * c
		d += total
		if d > float64(depth) {
			return pf, nil
		}
	}
	return 0.0, errors.New("has no enough depth")
}

func (p *Poloniex) limitBuy(amount, price string, currency string) (*proto.Order, error) {
	return p.placeLimitOrder(BUY, amount, price, currency)
}

func (p *Poloniex) limitSell(amount, price string, currency string) (*proto.Order, error) {
	return p.placeLimitOrder(SELL, amount, price, currency)
}

func (p *Poloniex) placeLimitOrder(command Command, amount, price string, currency string) (*proto.Order, error) {
	v := url.Values{}
	v.Set("command", command.String())
	v.Set("currencyPair", strings.ToUpper(currency))
	v.Set("rate", price)
	v.Set("amount", amount)
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	sign, err := p.buildPostForm(&v)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("Key", p.accessKey)
	header.Add("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", tradeAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}
	tresp := new(TradeResponse)
	err = json.Unmarshal(resp, &tresp)
	if err != nil {
		return nil, err
	}

	return &proto.Order{
		OrderID:      tresp.OrderNumber,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Price:        float64(0),
		Amount:       float64(0),
		DealedAmount: float64(0),
		Fee:          float64(0),
		Status:       proto.ORDER_UNFINISH,
		Currency:     currency,
		Side:         command.String(),
	}, nil
}

func (p *Poloniex) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return p.limitBuy(amount, price, currencyPair)
}

func (p *Poloniex) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return p.limitSell(amount, price, currencyPair)
}

func (p *Poloniex) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	v := url.Values{}
	v.Set("command", ORDERTRADES)
	v.Set("orderNumber", orderId)
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	sign, _ := p.buildPostForm(&v)

	header := http.Header{}
	header.Add("Key", p.accessKey)
	header.Add("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", tradeAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(resp), "error") {
		return &proto.Order{}, nil
	}
	or := new(OrderResponse)
	err = json.Unmarshal(resp, or)
	if err != nil {
		return nil, err
	}

	order := new(proto.Order)
	order.OrderID = orderId
	order.Currency = currencyPair
	for _, _ = range []TradeTerm(*or) {

	}

	return &proto.Order{}, nil
}

func (p *Poloniex) GetUnfinishOrders(currency string) ([]OpenOrder, error) {
	v := url.Values{}
	v.Set("command", OPENORDERS)
	v.Set("currencyPair", currency)
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	sign, err := p.buildPostForm(&v)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("Key", p.accessKey)
	header.Add("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", tradeAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}

	oos := make([]OpenOrder, 1)
	err = json.Unmarshal(resp, oos)
	if err != nil {
		return nil, err
	}

	return oos, nil
}

func (p *Poloniex) CancelOrder(orderId, currencypair string) (bool, error) {
	v := url.Values{}
	v.Set("command", CANCLEORDER)
	v.Set("orderNumber", orderId)
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	sign, err := p.buildPostForm(&v)
	if err != nil {
		return false, err
	}

	header := http.Header{}
	header.Add("Key", p.accessKey)
	header.Add("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", tradeAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
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
	return p.getAccount()
}

func (p *Poloniex) getAccount() (*proto.Account, error) {
	v := url.Values{}
	v.Add("command", COMPLETEBALANCES)
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	sign, err := p.buildPostForm(&v)
	if err != nil {
		return nil, err
	}

	header := http.Header{}
	header.Add("Key", p.accessKey)
	header.Add("Sign", sign)

	body := strings.NewReader(v.Encode())
	resp, err := util.Request("post", tradeAPI, "application/x-www-form-urlencoded", body, header, p.timeout)
	if err != nil {
		return nil, err
	}

	accountResp := make(map[string]*SubAccount)
	err = json.Unmarshal(resp, &accountResp)
	if err != nil || accountResp["error"] != nil {
		return nil, err
	}

	account := new(proto.Account)
	account.SubAccounts = make(map[string]proto.SubAccount)
	for k, v := range accountResp {
		account.Asset = 0.0
		account.Bourse = proto.Poloniex
		var subAccount proto.SubAccount
		avai, _ := strconv.ParseFloat(v.Available, 32)
		subAccount.Available = avai
		forz, _ := strconv.ParseFloat(v.OnOrders, 32)
		subAccount.Forzen = forz
		subAccount.Currency = k
		account.SubAccounts[strings.ToLower(k)] = subAccount
	}
	return account, nil
}

func (p *Poloniex) buildPostForm(v *url.Values) (string, error) {
	payload := v.Encode()
	sign, err := util.SHA512Sign(p.secretkey, payload)
	if err != nil {
		return "", err
	}
	return sign, nil
}
