package huobiN

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

var (
	API_URL          = "be.huobi.com"
	TICKER_URL       = "/market/kline?symbol=%s&period=1min"
	DEPTH_URL        = "/market/depth?symbol=%s&type=step1"
	GET_ACCOUNT_API  = "/v1/account/accounts/%s/balance"
	CREATE_ORDER_API = "/v1/order/orders"
	PLACE_ORDER_API  = "/v1/order/orders/%d/place"
	GET_ORDER_API    = "/v1/order/orders/%s"
	CANCEL_ORDER_API = "/v1/order/orders/%s/submitcancel"
)

type Huobi struct {
	accountid string
	accessKey string
	secretKey string
	timeout   int
}

func NewHuobi(accountid, accessKey, secretKey string, timeout int) (*Huobi, error) {
	return &Huobi{
		accountid: accountid,
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (huobi *Huobi) GetTicker(currencyPair string) (float64, error) {
	return 0, nil
}

func (huobi *Huobi) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf("https://%s%s", API_URL, fmt.Sprintf(DEPTH_URL, huobi.convertCurrencyPair(currencyPair)))
	rep, err := util.Request("GET", url, "application/json", nil, nil, huobi.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.HuobiN, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.HuobiN, currencyPair, err)
	}
	if body.Status != "ok" {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.HuobiN, currencyPair, body.Status)
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
			Sell:    sellprice,
			Buy:     buyprice,
			Sellnum: sellsum,
			Buynum:  buysum,
		}, nil
	}
	return nil, fmt.Errorf("sum not enough %v %v", sellsum, depth)
}

func (huobi *Huobi) GetAccount() (*proto.Account, error) {
	path := fmt.Sprintf(GET_ACCOUNT_API, huobi.accountid)
	params := url.Values{}
	sign, err := huobi.buildPostForm("GET", path, &params)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://%s%s?%s&Signature=%s", API_URL, path, params.Encode(), sign)
	rep, err := util.Request("GET", url, "application/json", nil, nil, huobi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetAccount err %v", err)
	}
	myaccount := MyAccount{}
	if err = json.Unmarshal(rep, &myaccount); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v", err)
	}
	if myaccount.Status != "ok" {
		return nil, fmt.Errorf("%s json Unmarshal err %s %s", proto.HuobiN, myaccount.Err_code, myaccount.Err_msg)
	}

	account := proto.Account{}
	//account.Asset =
	account.Bourse = strings.ToLower(proto.HuobiN)
	account.SubAccounts = make(map[string]proto.SubAccount)
	for _, a := range myaccount.Data.List {
		subAcc := proto.SubAccount{}
		if _, ok := account.SubAccounts[a.Currency]; ok {
			subAcc = account.SubAccounts[a.Currency]
		}
		if a.Type == "trade" {
			subAcc.Available, _ = strconv.ParseFloat(a.Balance, 64)
		} else if a.Type == "frozen" {
			subAcc.Forzen, _ = strconv.ParseFloat(a.Balance, 64)
		}
		subAcc.Currency = a.Currency
		account.SubAccounts[a.Currency] = subAcc
	}
	return &account, nil
}

func (huobi *Huobi) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	id, err := huobi.createOrder(side, amount, price, currencyPair)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	placeapi := fmt.Sprintf(PLACE_ORDER_API, id)
	sign, err := huobi.buildPostForm("POST", placeapi, &params)
	if err != nil {
		return nil, err
	}
	urls := fmt.Sprintf("https://%s%s?%s&Signature=%s", API_URL, placeapi, params.Encode(), sign)
	rep, err := util.Request("POST", urls, "application/json", nil, nil, huobi.timeout)
	if err != nil {
		return nil, fmt.Errorf("placeorder %s err:%v", side, err)
	}
	placeOrder := StatusOrder{}
	if err := json.Unmarshal(rep, &placeOrder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if placeOrder.Status != "ok" {
		return nil, fmt.Errorf("%s err:%v,%v,%v", side, placeOrder.Data, placeOrder.Err_code, placeOrder.Err_msg)
	}
	order := proto.Order{
		OrderID:  placeOrder.Data.(string),
		Currency: currencyPair,
	}
	return &order, nil
}

func (huobi *Huobi) createOrder(side, amount, price, currencyPair string) (int64, error) {
	var witchside string
	if side == proto.BUY {
		witchside = "buy-limit"
	} else if side == proto.SELL {
		witchside = "sell-limit"
	}
	signParams := url.Values{}
	sign, err := huobi.buildPostForm("POST", CREATE_ORDER_API, &signParams)
	if err != nil {
		return 0, err
	}
	urls := fmt.Sprintf("https://%s%s?%s&Signature=%s", API_URL, CREATE_ORDER_API, signParams.Encode(), sign)
	//create order
	params := CreateOrder{
		AccountId: huobi.accountid,
		Amount:    amount,
		Price:     price,
		Source:    "api",
		Symbol:    huobi.convertCurrencyPair(currencyPair),
		Type:      witchside,
	}
	body, _ := json.Marshal(params)
	rep, err := util.Request("POST", urls, "application/json", bytes.NewReader(body), nil, huobi.timeout)
	if err != nil {
		return 0, fmt.Errorf("createorder %s err:%v", side, err)
	}
	placeOrder := StatusOrder{}
	if err := json.Unmarshal(rep, &placeOrder); err != nil {
		return 0, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if placeOrder.Status != "ok" {
		return 0, fmt.Errorf("%s err:%v,%v,%v", side, placeOrder.Err_code, placeOrder.Err_msg, placeOrder.Data)
	}
	return int64(placeOrder.Data.(float64)), nil
}

func (huobi *Huobi) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return huobi.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (huobi *Huobi) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return huobi.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (huobi *Huobi) buildPostForm(method, path string, postForm *url.Values) (string, error) {
	postForm.Set("AccessKeyId", huobi.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format(proto.UTCTime))
	signBody := fmt.Sprintf("%s\n%s\n%s\n%s", method, API_URL, path, postForm.Encode())
	sign, err := util.SHA256SignByte(huobi.secretKey, signBody)
	if err != nil {
		return "", err
	}
	return url.QueryEscape(base64.StdEncoding.EncodeToString(sign)), nil
}

func (huobi *Huobi) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	getorder := fmt.Sprintf(GET_ORDER_API, orderId)
	sign, err := huobi.buildPostForm("GET", getorder, &params)
	if err != nil {
		return nil, err
	}
	urls := fmt.Sprintf("https://%s%s?%s&Signature=%s", API_URL, getorder, params.Encode(), sign)

	rep, err := util.Request("GET", urls, "application/json", nil, nil, huobi.timeout)
	if err != nil {
		return nil, fmt.Errorf("getorder %s err:%v", orderId, err)
	}
	myOrder := MyOrder{}
	if err := json.Unmarshal(rep, &myOrder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if myOrder.Status != "ok" {
		return nil, fmt.Errorf("getorder %s err:%v %v", orderId, myOrder.Err_code, myOrder.Err_msg)
	}

	Price, _ := strconv.ParseFloat(myOrder.Data.Price, 64)
	Amount, _ := strconv.ParseFloat(myOrder.Data.Amount, 64)
	DealedAmount, _ := strconv.ParseFloat(myOrder.Data.Field_amount, 64)
	Fee, _ := strconv.ParseFloat(myOrder.Data.Field_fees, 64)
	OrderTime := time.Unix(myOrder.Data.Created_at/1000, 0).Format(proto.LocalTime)
	var Status string
	switch myOrder.Data.State {
	case "pre-submitted", "submitting", "submitted": //准备提交,已提交
		Status = proto.ORDER_UNFINISH
	case "partial-filled": //部分成交
		Status = proto.ORDER_PART_FINISH
	case "filled": //完全成交
		Status = proto.ORDER_FINISH
	case "partial-canceled": //部分成交撤销
		Status = proto.ORDER_CANCEL_ING
	case "canceled": //已撤销
		Status = proto.ORDER_CANCEL
	}
	return &proto.Order{
		OrderID:      fmt.Sprintf("%d", myOrder.Data.ID),
		OrderTime:    OrderTime,
		Price:        Price,
		Amount:       Amount,
		DealedAmount: DealedAmount,
		Fee:          Fee,
		Status:       Status,
		Currency:     myOrder.Data.Symbol,
		Side:         myOrder.Data.Type,
	}, nil
}

func (huobi *Huobi) CancelOrder(orderId, currencyPair string) (bool, error) {
	params := url.Values{}
	cancelorder := fmt.Sprintf(CANCEL_ORDER_API, orderId)
	sign, err := huobi.buildPostForm("POST", cancelorder, &params)
	if err != nil {
		return false, err
	}
	urls := fmt.Sprintf("https://%s%s?%s&Signature=%s", API_URL, cancelorder, params.Encode(), sign)
	rep, err := util.Request("POST", urls, "application/json", nil, nil, huobi.timeout)
	if err != nil {
		return false, fmt.Errorf("cancel %s err:%v", orderId, err)
	}
	myOrder := StatusOrder{}
	if err := json.Unmarshal(rep, &myOrder); err != nil {
		return false, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if myOrder.Status != "ok" {
		return false, fmt.Errorf("cancel %s err:%v %v", orderId, myOrder.Err_code, myOrder.Err_msg)
	}
	return true, nil
}

func (huobi *Huobi) convertCurrencyPair(currencyPair string) string {
	switch currencyPair {
	case proto.ETH_CNY:
		return "ethcny"
	case proto.ETC_CNY:
		return "etccny"
	default:
		return ""
	}
}
