package bter

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

var (
	API_URL      = "https://data.bter.com/api2/1"
	TICKER_URL   = "/tickers"
	DEPTH_URL    = "/orderBook/%s"
	ORDER_BUY    = "/private/buy"
	ORDER_SELL   = "/private/sell"
	ORDER_GET    = "/private/getOrder"
	GET_ACCOUNT  = "/private/balances"
	CANCEL_ORDER = "/private/cancelOrder"
)

type Bter struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewBter(accessKey, secretKey string, timeout int) (*Bter, error) {
	return &Bter{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (bter *Bter) GetTicker(currencyPair string) (float64, error) {
	url := fmt.Sprintf(API_URL+TICKER_URL, currencyPair)
	rep, err := util.Request("GET", url, "application/json", nil, nil, bter.timeout)
	if err != nil {
		return 0, fmt.Errorf("%s request err %s %v", proto.Bter, currencyPair, err)
	}
	body := Market{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return 0, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bter, currencyPair, err)
	}
	return body.Last, nil
}

func (bter *Bter) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf(API_URL+DEPTH_URL, currencyPair)
	rep, err := util.Request("GET", url, "application/x-www-form-urlencoded", nil, nil, bter.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Bter, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bter, currencyPair, err)
	}
	if body.Result != "true" || len(body.Asks) == 0 || len(body.Bids) == 0 {
		return nil, fmt.Errorf("resault not true %s", body.Message)
	}

	var sellsum float64
	var sellprice float64
	var buysum float64
	var buyprice float64
	var len int = len(body.Asks)
	for i := len - 1; i >= 0; i-- {
		var sum float64
		switch param := (body.Asks[i][1]).(type) {
		case float64:
			sum = param
		case string:
			if sum, err = strconv.ParseFloat(param, 64); err != nil {
				continue
			}
		default:
			continue
		}
		sellsum += sum
		if sellsum > float64(depth) {
			switch param := (body.Asks[i][0]).(type) {
			case float64:
				sellprice = param
			case string:
				if sellprice, err = strconv.ParseFloat(param, 64); err != nil {
					continue
				}
			default:
				continue
			}
			break
		}
	}

	for i := 0; i < len; i++ {
		var sum float64
		switch param := (body.Bids[i][1]).(type) {
		case float64:
			sum = param
		case string:
			if sum, err = strconv.ParseFloat(param, 64); err != nil {
				continue
			}
		default:
			continue
		}
		buysum += sum
		if buysum > float64(depth) {
			switch param := (body.Bids[i][0]).(type) {
			case float64:
				buyprice = param
			case string:
				if buyprice, err = strconv.ParseFloat(param, 64); err != nil {
					continue
				}
			default:
				continue
			}
			break
		}
	}
	if sellsum > float64(depth) && buysum > float64(depth) {
		return &proto.Price{
			Sell:    sellprice,
			Buy:     buyprice,
			Sellnum: sellsum,
			Buynum:  buysum,
		}, nil
	}
	return nil, fmt.Errorf("%s sum not enough %v %v", proto.Bter, sellsum, depth)
}

func (bter *Bter) GetAccount() (*proto.Account, error) {
	var params url.Values
	Sign, err := bter.buildPostForm(&params)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Key", bter.accessKey)
	header.Set("Sign", Sign)

	rep, err := util.Request("POST", API_URL+GET_ACCOUNT,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		header, bter.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetAccount err %v", err)
	}
	myaccount := MyAccount{}
	err = json.Unmarshal(rep, &myaccount)
	if err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v", err)
	}
	if myaccount.Result != "true" {
		return nil, fmt.Errorf("account false err %v", myaccount.Result)
	}

	account := proto.Account{}
	//account.Asset
	account.Bourse = strings.ToLower(proto.Bter)
	account.SubAccounts = make(map[string]proto.SubAccount)

	//cny
	subAcc := proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.CNY, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.CNY, 64)
	subAcc.Currency = proto.CNY
	account.SubAccounts[subAcc.Currency] = subAcc
	//btc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.BTC, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.BTC, 64)
	subAcc.Currency = proto.BTC
	account.SubAccounts[subAcc.Currency] = subAcc
	//etc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.ETC, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.ETC, 64)
	subAcc.Currency = proto.ETC
	account.SubAccounts[subAcc.Currency] = subAcc
	//eth
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.ETH, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.ETH, 64)
	subAcc.Currency = proto.ETH
	account.SubAccounts[subAcc.Currency] = subAcc
	//ltc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.LTC, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.LTC, 64)
	subAcc.Currency = proto.LTC
	account.SubAccounts[subAcc.Currency] = subAcc
	//snt
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.SNT, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.SNT, 64)
	subAcc.Currency = proto.SNT
	account.SubAccounts[subAcc.Currency] = subAcc
	//omg
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.OMG, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.OMG, 64)
	subAcc.Currency = proto.OMG
	account.SubAccounts[subAcc.Currency] = subAcc
	//pay
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.PAY, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.PAY, 64)
	subAcc.Currency = proto.PAY
	account.SubAccounts[subAcc.Currency] = subAcc
	//btm
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.Available.BTM, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.Locked.BTM, 64)
	subAcc.Currency = proto.BTM
	account.SubAccounts[subAcc.Currency] = subAcc
	return &account, nil
}

func (bter *Bter) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("currencyPair", currencyPair)
	params.Set("rate", price)
	params.Set("amount", amount)

	Sign, err := bter.buildPostForm(&params)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Key", bter.accessKey)
	header.Set("Sign", Sign)

	var SIDEURL string
	if side == proto.SELL {
		SIDEURL = ORDER_SELL
	} else if side == proto.BUY {
		SIDEURL = ORDER_BUY
	}

	rep, err := util.Request("POST", API_URL+SIDEURL,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		header, bter.timeout)
	if err != nil {
		return nil, fmt.Errorf("request %s err:%v", side, err)
	}
	placeOrder := PlaceOrder{}
	if err := json.Unmarshal(rep, &placeOrder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if placeOrder.Result != "true" {
		return nil, fmt.Errorf("%s err:%v", side, placeOrder.Message)
	}
	myorder := MyOrder{}
	myorder.Order.OrderNumber = fmt.Sprintf("%d", placeOrder.OrderNumber)
	myorder.Order.CurrencyPair = currencyPair
	return bter.parseOrder(&myorder)
}

func (bter *Bter) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return bter.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (bter *Bter) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return bter.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (bter *Bter) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("orderNumber", orderId)
	params.Set("currencyPair", currencyPair)
	Sign, err := bter.buildPostForm(&params)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Key", bter.accessKey)
	header.Set("Sign", Sign)

	rep, err := util.Request("POST", API_URL+ORDER_GET,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		header, bter.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetOneOrder err %v %s", err, orderId)
	}
	myorder := MyOrder{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if myorder.Result != "true" {
		return nil, fmt.Errorf("GetOneOrder err id:%s %s time:%s", orderId, myorder.Message, myorder.Elapsed)
	}
	return bter.parseOrder(&myorder)
}

func (bter *Bter) CancelOrder(orderId, currencyPair string) (bool, error) {
	params := url.Values{}
	params.Set("orderNumber", orderId)
	params.Set("currencyPair", currencyPair)
	Sign, err := bter.buildPostForm(&params)
	if err != nil {
		return false, err
	}
	header := http.Header{}
	header.Set("Key", bter.accessKey)
	header.Set("Sign", Sign)

	rep, err := util.Request("POST", API_URL+CANCEL_ORDER,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		header, bter.timeout)
	if err != nil {
		return false, fmt.Errorf("request CancelOrder err %s", err)
	}
	if rep == nil {
		return false, err
	}
	cancel := CancelOrder{}
	if err := json.Unmarshal(rep, &cancel); err != nil {
		return false, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if cancel.Code != 0 {
		return false, fmt.Errorf("CancelOrder err id:%s %s", orderId, cancel.Message)
	}
	return true, nil
}

func (bter *Bter) buildPostForm(postForm *url.Values) (string, error) {
	payload := postForm.Encode()
	return util.SHA512Sign(bter.secretKey, payload)
}

func (bter *Bter) parseOrder(myorder *MyOrder) (*proto.Order, error) {
	var status string
	var Price float64
	switch myorder.Order.Status {
	case "open":
		status = proto.ORDER_UNFINISH
	case "closed":
		status = proto.ORDER_FINISH
	default:
		status = myorder.Order.Status
	}
	if status == proto.ORDER_FINISH {
		Price, _ = strconv.ParseFloat(myorder.Order.FilledRate.(string), 64)
	} else {
		Price = myorder.Order.InitialRate
	}
	Amount, _ := strconv.ParseFloat(myorder.Order.InitialAmount, 64)
	return &proto.Order{
		Amount:       Amount * 1.5,
		OrderID:      myorder.Order.OrderNumber,
		Price:        Price,
		DealedAmount: myorder.Order.FilledAmount * 1.5,
		//Fee:          myorder.Order.,
		Currency:  myorder.Order.CurrencyPair,
		Status:    status,
		OrderTime: time.Now().Format(proto.LocalTime),
		Side:      myorder.Order.Type,
	}, nil
}
