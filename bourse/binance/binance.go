package binance

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"errors"
	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
	"net/http"
)

const (
	BASE_URL    = "https://api.binance.com/"
	TICKER_API  = ""
	DEPTH_API   = "api/v1/depth"
	ACCOUNT_API = "api/v3/account"
	ORDER_API   = "api/v3/order"
)

type Binance struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewBinance(accessKey, secretKey string, timeout int) (*Binance, error) {
	return &Binance{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (binance *Binance) GetTicker(currencyPair string) (float64, error) {
	url := fmt.Sprintf(BASE_URL + fmt.Sprintf(TICKER_API, convertCurrency(currencyPair)))
	rep, err := util.Request("GET", url, "application/json", nil, nil, binance.timeout)
	if err != nil {
		return 0, fmt.Errorf("%s request err %s %v", proto.Binance, currencyPair, err)
	}
	body := Market{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return 0, fmt.Errorf("%s json Unmarshal err %s %v", proto.Binance, currencyPair, err)
	}
	return strconv.ParseFloat(body.Ticker.Last, 64)
}

func (binance *Binance) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf("%s/%s?symbol=%s", BASE_URL, DEPTH_API, convertCurrency(currencyPair))
	resp, err := util.Request("GET", url, "application/json", nil, nil, 4)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Binance, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Binance, currencyPair, err)
	}
	if body.Error != "" {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Binance, currencyPair, body.Error)
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
		c := entry[1].(string)
		price, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0.0, 0.0, err
		}
		amount, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return 0.0, 0.0, err
		}
		total := price * amount
		sum += total
		if sum > float64(depth) {
			return price, sum, nil
		}
	}
	return 0.0, 0.0, fmt.Errorf("has no enough depth sum:%v depth:%v", sum, depth)
}

func (binance *Binance) GetAccount() (*proto.Account, error) {
	h := http.Header{}
	h.Add("X-MBX-APIKEY", binance.accessKey)
	//symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559
	params := url.Values{}
	binance.buildParamsSigned(&params)
	url := fmt.Sprintf("%s%s?%s", BASE_URL, ACCOUNT_API, params.Encode())
	b, err := util.HttpGet(url, &h)
	if err != nil {
		return nil, err
	}

	fmt.Println("account resp:", string(b))

	var dataMap map[string]interface{}
	err = json.Unmarshal(b, &dataMap)
	if err != nil {
		return nil, err
	}

	if _, isok := dataMap["code"]; isok == true {
		return nil, errors.New(dataMap["msg"].(string))
	}
	account := proto.Account{}
	account.SubAccounts = make(map[string]proto.SubAccount)

	balances := dataMap["balances"].([]interface{})
	for _, v := range balances {
		//log.Println(v)
		vv := v.(map[string]interface{})
		currency := vv["asset"].(string)
		account.SubAccounts[currency] = proto.SubAccount{
			Currency:  currency,
			Available: util.ToFloat64(vv["free"]),
			Forzen:    util.ToFloat64(vv["locked"]),
		}
	}

	return &account, nil
}

func (binance *Binance) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	fmt.Println(strings.ToUpper(fmt.Sprintf("%s", side)))
	params := url.Values{}
	params.Set("symbol", convertCurrency(currencyPair))
	params.Set("side", strings.ToUpper(fmt.Sprintf("%s", side)))
	params.Set("price", price)
	params.Set("quantity", amount)
	params.Set("type", "LIMIT")
	params.Set("timeInForce", "GTC")
	binance.buildParamsSigned(&params)

	h := http.Header{}
	h.Add("X-MBX-APIKEY", binance.accessKey)
	//url := fmt.Sprintf("%s%s?")
	resp, err := util.Request("POST", BASE_URL+ORDER_API,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		h, binance.timeout)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp))
	placeOrder := PlaceOrder{}
	if err := json.Unmarshal(resp, &placeOrder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(resp))
	}
	// !=0说明有错误
	if placeOrder.Code != 0 {
		return nil, fmt.Errorf("%s err:%v", placeOrder.Code, placeOrder.Msg)
	}
	return &proto.Order{
		OrderID:  fmt.Sprintf("%d", placeOrder.OrderId),
		Currency: placeOrder.Symbol,
	}, nil
}

func (binance *Binance) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return binance.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (binance *Binance) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return binance.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (binance *Binance) CancelOrder(orderId, currencyPair string) (bool, error) {
	params := url.Values{}
	params.Set("symbol", convertCurrency(currencyPair))
	params.Set("orderId", orderId)

	binance.buildParamsSigned(&params)
	h := http.Header{}
	h.Add("X-MBX-APIKEY", binance.accessKey)
	h.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := util.Request("DELETE", BASE_URL+ORDER_API, "", strings.NewReader(params.Encode()), h, binance.timeout)
	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	orderIdCanceled := util.ToInt(respmap["orderId"])
	if orderIdCanceled <= 0 {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (binance *Binance) parseOrder(myorder *MyOrder) (*proto.Order, error) {
	var status string
	switch myorder.Status {
	case 0:
		status = proto.ORDER_UNFINISH
	case 1:
		status = proto.ORDER_CANCEL
	case 2:
		status = proto.ORDER_FINISH
	case 3:
		status = proto.ORDER_PART_FINISH
	}
	var Side string
	if myorder.Type == 1 {
		Side = proto.BUY
	} else {
		Side = proto.SELL
	}
	return &proto.Order{
		Amount:       myorder.TradeAmount,
		Fee:          myorder.Fees,
		OrderID:      myorder.ID,
		Price:        float64(myorder.Price),
		DealedAmount: myorder.TradeAmount,
		Currency:     myorder.Currency,
		Status:       status,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Side:         Side,
	}, nil
}

func (binance *Binance) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("method", "getOrder")
	params.Set("id", orderId)
	params.Set("currency", currencyPair)
	binance.buildParamsSigned(&params)

	rep, err := util.Request("POST", "",
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, binance.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetOneOrder err %v %s", err, orderId)
	}
	myorder := MyOrder{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if myorder.Code > 0 {
		return nil, fmt.Errorf("GetOneOrder err id:%s %s", orderId, myorder.Message)
	}
	return binance.parseOrder(&myorder)
}

func Deposit(currency string, amount float64) error {
	return nil
}
func Withdraw(currency string, amount float64) error {
	return nil
}

func (binance *Binance) buildParamsSigned(postForm *url.Values) error {
	postForm.Set("recvWindow", "6000000")
	tonce := strconv.FormatInt(time.Now().UnixNano(), 10)[0:13]
	postForm.Set("timestamp", tonce)
	payload := postForm.Encode()
	sign, _ := util.SHA256Sign(binance.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}

func convertCurrency(currency string) string {
	return strings.Trim(currency,"_")
}

