package yunbi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

const _EXCHANGE_NAME = "yunbi.com"

var (
	API_URL          = "https://yunbi.com"
	API_URI_PREFIX   = "/api/v2/"
	TICKER_URL       = "tickers/%s.json"
	DEPTH_URL        = "depth.json?market=%s&limit=%d"
	USER_INFO_URL    = "members/me.json"
	GET_ORDER_API    = "order.json"
	DELETE_ORDER_API = "order/delete.json"
	PLACE_ORDER_API  = "orders.json"
)

type YunBi struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewYunBi(accessKey, secretKey string, timeout int) (*YunBi, error) {
	return &YunBi{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (yunbi *YunBi) GetTicker(currencyPair string) (float64, error) {
	url := fmt.Sprintf(API_URL+API_URI_PREFIX+TICKER_URL, yunbi.convertCurrencyPair(currencyPair))
	rep, err := util.Request("GET", url, "application/json", nil, nil, yunbi.timeout)
	if err != nil {
		return 0, fmt.Errorf("%s request err %s %v", proto.Yunbi, currencyPair, err)
	}
	body := Market{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return 0, fmt.Errorf("%s json Unmarshal err %s %v", proto.Yunbi, currencyPair, err)
	}
	return body.Ticker.Last, nil
}

func (yunbi *YunBi) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf(API_URL+API_URI_PREFIX+DEPTH_URL, yunbi.convertCurrencyPair(currencyPair), size)
	rep, err := util.Request("GET", url, "application/json", nil, nil, yunbi.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Yunbi, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Yunbi, currencyPair, err)
	}

	var sellsum float64
	var buysum float64
	var sellprice float64
	var buyprice float64

	var len int = len(body.Asks)
	for i := len - 1; i >= 0; i-- {
		price, err := strconv.ParseFloat((body.Asks[i][0]), 64)
		if err != nil {
			continue
		}
		sum, err := strconv.ParseFloat((body.Asks[i][1]), 64)
		if err != nil {
			continue
		}
		sellsum += sum
		if sellsum > float64(depth) {
			sellprice = price
			break
		}
	}

	for i := 0; i < len; i++ {
		price, err := strconv.ParseFloat((body.Bids[i][0]), 64)
		if err != nil {
			continue
		}
		sum, err := strconv.ParseFloat((body.Bids[i][1]), 64)
		if err != nil {
			continue
		}
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
		//data, _ := json.Marshal(price)
		//log.Debug("yunbi body ", string(data))
		//return fmt.Sprintf("etc\nselldepth %0.3f price:%0.3f\nbuydepth %0.3f price:%0.3f", sellsum, sellprice, buysum, buyprice)
		return &price, nil
	}
	return nil, fmt.Errorf("sum not enough sell:%v buy:%v depth:%v", sellsum, buysum, depth)
}

func (yunbi *YunBi) GetAccount() (*proto.Account, error) {
	urls := API_URL + API_URI_PREFIX + USER_INFO_URL
	params := url.Values{}
	yunbi.buildPostForm("GET", API_URI_PREFIX+USER_INFO_URL, &params)

	rep, err := util.Request("GET", urls+"?"+params.Encode(), "application/json", nil, nil, yunbi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetAccount err %v", err)
	}
	myaccount := MyAccount{}
	err = json.Unmarshal(rep, &myaccount)
	if err != nil {
		return nil, fmt.Errorf("json Unmarshal err %s", err)
	}

	account := proto.Account{}
	account.Bourse = strings.ToLower(proto.Yunbi)
	//account.Asset = 0.0 //需要计算
	account.SubAccounts = make(map[string]proto.SubAccount)
	for _, a := range myaccount.Accounts {
		subAcc := proto.SubAccount{}
		subAcc.Available, _ = strconv.ParseFloat(a.Balance, 64)
		subAcc.Forzen, _ = strconv.ParseFloat(a.Locked, 64)
		subAcc.Currency = a.Currency
		account.SubAccounts[a.Currency] = subAcc
	}
	return &account, nil
}

func (yunbi *YunBi) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("market", yunbi.convertCurrencyPair(currencyPair))
	params.Set("side", side)
	params.Set("price", price)
	params.Set("volume", amount)
	yunbi.buildPostForm("POST", API_URI_PREFIX+PLACE_ORDER_API, &params)
	rep, err := util.Request("POST", API_URL+API_URI_PREFIX+PLACE_ORDER_API,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, yunbi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request %s err:%v", side, err)
	}
	//log.Debugf("%s %s", side, string(rep))
	myorder := MyOrder{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	return yunbi.parseOrder(&myorder)
}

func (yunbi *YunBi) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return yunbi.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (yunbi *YunBi) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return yunbi.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (yunbi *YunBi) CancelOrder(orderId, currencyPair string) (bool, error) {
	params := url.Values{}
	params.Set("id", orderId)
	yunbi.buildPostForm("POST", API_URI_PREFIX+DELETE_ORDER_API, &params)

	rep, err := util.Request("POST", API_URL+API_URI_PREFIX+DELETE_ORDER_API,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, yunbi.timeout)
	if err != nil {
		return false, fmt.Errorf("request CancelOrder err %s", err)
	}
	if rep == nil {
		return false, err
	}
	//log.Debug("cancel:", string(rep))
	return true, nil
}

func (yunbi *YunBi) parseOrder(myorder *MyOrder) (*proto.Order, error) {
	Price, _ := strconv.ParseFloat(myorder.Price, 64)
	DealedAmount, _ := strconv.ParseFloat(myorder.ExecutedVolume, 64)
	amount, _ := strconv.ParseFloat(myorder.Volume, 64)
	return &proto.Order{
		OrderID:      fmt.Sprintf("%d", myorder.ID),
		Price:        Price,
		Amount:       amount,
		Currency:     myorder.Market,
		DealedAmount: DealedAmount,
		Status:       myorder.State,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Side:         myorder.Side,
	}, nil
	//order.Fee = 计算
	//log.Debug("order price:", order.Price, "send price:", price) //对比执行完订单和下发的区别
}

func (yunbi *YunBi) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("id", orderId)
	yunbi.buildPostForm("GET", API_URI_PREFIX+GET_ORDER_API, &params)

	rep, err := util.Request("GET", API_URL+API_URI_PREFIX+GET_ORDER_API+"?"+params.Encode(),
		"application/x-www-form-urlencoded", nil, nil, yunbi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetOneOrder err %v", err)
	}
	//log.Debug("getorder:", string(rep))

	myorder := MyOrder{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	return yunbi.parseOrder(&myorder)
}

func (yunbi *YunBi) GetUnfinishOrders(currencyPair string) (*[]proto.Order, error) {
	params := url.Values{}
	params.Set("market", yunbi.convertCurrencyPair(currencyPair))
	params.Set("state", "wait")
	yunbi.buildPostForm("GET", API_URI_PREFIX+PLACE_ORDER_API, &params)

	rep, err := util.Request("GET", API_URL+API_URI_PREFIX+PLACE_ORDER_API+"?"+params.Encode(),
		"application/x-www-form-urlencoded", nil, nil, yunbi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetUnfinishOrders err %v", err)
	}

	myorders := []MyOrder{}
	if err := json.Unmarshal(rep, &myorders); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	orders := []proto.Order{}
	for _, myorder := range myorders {
		if order, err := yunbi.parseOrder(&myorder); err != nil {
			orders = append(orders, *order)
		}
	}
	return &orders, nil
}

func (yunbi *YunBi) buildPostForm(httpMethod, apiURI string, postForm *url.Values) error {
	postForm.Set("access_key", yunbi.accessKey)
	postForm.Set("tonce", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))

	params := postForm.Encode()
	payload := httpMethod + "|" + apiURI + "|" + params
	//println(payload)

	sign, err := util.SHA256Sign(yunbi.secretKey, payload)
	if err != nil {
		return err
	}

	postForm.Set("signature", sign)
	//postForm.Del("secret_key")
	return nil
}

func (y *YunBi) convertCurrencyPair(currencyPair string) string {
	switch currencyPair {
	case proto.BTC_CNY:
		return "btccny"
	case proto.ETH_CNY:
		return "ethcny"
	case proto.ETC_CNY:
		return "etccny"
	case proto.LTC_CNY:
		return "ltccny"
	case proto.EOS_CNY:
		return "eoscny"
	}
	return "btccny"
}
