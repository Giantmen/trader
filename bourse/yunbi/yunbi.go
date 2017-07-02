package yunbi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Giantmen/trader/config"
	"github.com/Giantmen/trader/log"
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

func NewYunBi(cfg *config.Server) (*YunBi, error) {
	return &YunBi{
		accessKey: cfg.Accesskey,
		secretKey: cfg.Secretkey,
		timeout:   cfg.Timeout,
	}, nil
}

func (yunbi *YunBi) GetTicker(currencyPair string) (float64, error) {
	url := fmt.Sprintf(API_URL+API_URI_PREFIX+TICKER_URL, yunbi.convertCurrencyPair(currencyPair))
	rep, err := util.Request("GET", url, "application/json", nil, nil, yunbi.timeout)
	if err != nil {
		log.Error("request err", proto.Yunbi, currencyPair, err)
		return 0, err
	}
	body := Market{}
	if err := json.Unmarshal(rep, &body); err != nil {
		log.Error("json Unmarshal err ", proto.Yunbi, currencyPair, err)
		return 0, err
	}
	return body.Ticker.Last, nil
}

func (yunbi *YunBi) GetPriceOfDepth(size, depth int, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf(API_URL+API_URI_PREFIX+DEPTH_URL, yunbi.convertCurrencyPair(currencyPair), size)
	rep, err := util.Request("GET", url, "application/json", nil, nil, yunbi.timeout)
	if err != nil {
		log.Error("request err", proto.Yunbi, currencyPair, err)
		return nil, err
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		log.Error("json Unmarshal err ", proto.Yunbi, currencyPair, err)
		return nil, err
	}

	var sellsum float64
	var buysum float64
	var sellprice float64
	var buyprice float64

	var len int = len(body.Asks)
	for i := len - 1; i >= 0; i-- {
		price, err := strconv.ParseFloat((body.Asks[i][0]), 64)
		if err != nil {
			log.Error("err", err)
			continue
		}
		sum, err := strconv.ParseFloat((body.Asks[i][1]), 64)
		if err != nil {
			log.Error("err", err)
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
			log.Error("err", err)
			continue
		}
		sum, err := strconv.ParseFloat((body.Bids[i][1]), 64)
		if err != nil {
			log.Error("err", err)
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
		data, _ := json.Marshal(price)
		log.Debug("yunbi body ", string(data))
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
		log.Error("request GetAccount err", err)
		return nil, err
	}
	myaccount := MyAccount{}
	err = json.Unmarshal(rep, &myaccount)
	if err != nil {
		log.Error("json Unmarshal err", err)
		return nil, err
	}

	account := proto.Account{}
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

func (yunbi *YunBi) placeOrder(side, amount, price string, currencyPair string) (*proto.Order, error) {
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
		log.Error("request GetAccount err", err)
		return nil, err
	}
	order, err := yunbi.parseOrder(rep)
	if err != nil {
		return nil, err
	}
	order.Currency = currencyPair
	order.Amount, _ = strconv.ParseFloat(amount, 64)
	//order.Fee = 计算
	//log.Debug("order price:", order.Price, "send price:", price) //对比执行完订单和下发的区别
	return order, nil
}

func (yunbi *YunBi) Buy(amount, price string, currencyPair string) (*proto.Order, error) {
	return yunbi.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (yunbi *YunBi) Sell(amount, price string, currencyPair string) (*proto.Order, error) {
	return yunbi.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (yunbi *YunBi) CancelOrder(orderId string, currencyPair string) (bool, error) {
	params := url.Values{}
	params.Set("id", orderId)
	yunbi.buildPostForm("POST", API_URI_PREFIX+DELETE_ORDER_API, &params)

	rep, err := util.Request("POST", API_URL+API_URI_PREFIX+DELETE_ORDER_API,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, yunbi.timeout)
	if err != nil {
		log.Error("request GetAccount err", err)
		return false, err
	}

	log.Debug("rep:", string(rep)) //////////////////////debug
	// body := xxxx{}
	// err = json.Unmarshal(rep, &respMap)
	// if err != nil {
	// 	log.Error("json Unmarshal err", err)
	// 	return false, err
	// }
	return true, nil
}
func (yunbi *YunBi) parseOrder(rep []byte) (*proto.Order, error) {
	myorder := MyOrder{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		log.Error("json Unmarshal err", err, string(rep))
		return nil, err
	}

	Price, _ := strconv.ParseFloat(myorder.Price, 64)
	DealedAmount, _ := strconv.ParseFloat(myorder.ExecutedVolume, 64)
	return &proto.Order{
		OrderID:      fmt.Sprintf("%d", myorder.ID),
		Price:        Price,
		DealedAmount: DealedAmount,
		Currency:     myorder.Market,
		Status:       myorder.State,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Side:         myorder.Side,
	}, nil
}

func (yunbi *YunBi) GetOneOrder(orderId string, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("id", orderId)
	yunbi.buildPostForm("GET", API_URI_PREFIX+GET_ORDER_API, &params)

	rep, err := util.Request("GET", API_URL+API_URI_PREFIX+GET_ORDER_API+"?"+params.Encode(),
		"application/x-www-form-urlencoded", nil, nil, yunbi.timeout)
	if err != nil {
		log.Error("request GetAccount err", err)
		return nil, err
	}

	log.Debug(string(rep))

	order, err := yunbi.parseOrder(rep)
	if err != nil {
		return nil, err
	}
	order.Currency = currencyPair
	return order, nil
}

func (yunbi *YunBi) GetUnfinishOrders(currencyPair string) (*[]proto.Order, error) {
	params := url.Values{}
	params.Set("market", yunbi.convertCurrencyPair(currencyPair))
	params.Set("state", "wait")
	yunbi.buildPostForm("GET", API_URI_PREFIX+PLACE_ORDER_API, &params)

	rep, err := util.Request("GET", API_URL+API_URI_PREFIX+PLACE_ORDER_API+"?"+params.Encode(),
		"application/x-www-form-urlencoded", nil, nil, yunbi.timeout)
	if err != nil {
		log.Error("request GetAccount err", err)
		return nil, err
	}

	myorders := []MyOrder{}
	if err := json.Unmarshal(rep, &myorders); err != nil {
		log.Error("json Unmarshal err", err, string(rep))
		return nil, err
	}
	orders := []proto.Order{}
	for _, myorder := range myorders {
		Price, _ := strconv.ParseFloat(myorder.Price, 64)
		DealedAmount, _ := strconv.ParseFloat(myorder.ExecutedVolume, 64)
		order := proto.Order{
			OrderID:      fmt.Sprintf("%d", myorder.ID),
			Price:        Price,
			DealedAmount: DealedAmount,
			Currency:     myorder.Market,
			Status:       myorder.State,
			OrderTime:    myorder.CreatedAt,
			Side:         myorder.Side,
		}
		orders = append(orders, order)
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
	}
	return "btccny"
}
