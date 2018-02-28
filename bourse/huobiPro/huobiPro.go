package huobiPro

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
	//"strconv"
)

const (
	BASE_URL = "https://api.huobi.pro"
)

type HuobiPro struct {
	accountID string
	accessKey string
	secretKey string
	timeout   int
}

func NewHuobiPro(accountID,accessKey, secretKey string, timeout int) (*HuobiPro, error) {
	return &HuobiPro{
		accountID:accountID,
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (hbp *HuobiPro) GetAccountId() (string, error) {
	path := "/v1/account/accounts"
	params := &url.Values{}
	hbp.buildPostForm("GET", path, params)

	//log.Println(BASE_URL + path + "?" + params.Encode())

	resp, err := util.HttpGet(BASE_URL+path+"?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return "", nil
	}
	if respmap["status"].(string) != "ok" {
		return "", errors.New(respmap["err-code"].(string))
	}

	data := respmap["data"].([]interface{})
	accountIdMap := data[0].(map[string]interface{})
	hbp.accountID = fmt.Sprintf("%.f", accountIdMap["id"].(float64))

	//log.Println(respmap)
	return hbp.accountID, nil
}

func (hbp *HuobiPro) GetAccount() (*proto.Account, error) {
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", hbp.accountID)
	//path := fmt.Sprintf("/v1/account/accounts")
	params := &url.Values{}
	params.Set("accountId-id", hbp.accountID)
	hbp.buildPostForm("GET", path, params)

	urlStr := BASE_URL + path + "?" + params.Encode()
	resp, err := util.HttpGet(urlStr, nil)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(resp))
	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}
	if respmap["status"].(string)== "error" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	if datamap["state"].(string) != "working" {
		return nil, errors.New(datamap["state"].(string))
	}

	list := datamap["list"].([]interface{})
	acc := new(proto.Account)
	acc.SubAccounts = make(map[string]proto.SubAccount, 6)
	//acc.Exchange = hbp.GetExchangeName()

	subAccMap := make(map[string]*proto.SubAccount)

	for _, v := range list {
		balancemap := v.(map[string]interface{})
		currency := balancemap["currency"].(string)

		typeStr := balancemap["type"].(string)
		balance := util.ToFloat64(balancemap["balance"])
		if subAccMap[currency] == nil {
			subAccMap[currency] = new(proto.SubAccount)
		}
		subAccMap[currency].Currency = currency
		switch typeStr {
		case "trade":
			subAccMap[currency].Available = balance
		case "frozen":
			subAccMap[currency].Forzen = balance
		}
	}

	for k, v := range subAccMap {
		acc.SubAccounts[k] = *v
	}

	return acc, nil
}

func (hbp *HuobiPro) placeOrder(side string, amount, price string, currencyPair string, orderType string) (*proto.Order, error) {
	path := "/v1/order/orders/place"
	params := url.Values{}
	params.Set("account-id", hbp.accountID)
	params.Set("amount", amount)
	params.Set("symbol", currencyPair)
	params.Set("type", orderType)

	switch orderType {
	case "buy-limit", "sell-limit":
		params.Set("price", price)
	}

	hbp.buildPostForm("POST", path, &params)

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Language", "zh-cn")

	resp, err := util.Request("POST", BASE_URL+path+"?"+params.Encode(), "application/json", strings.NewReader(hbp.toJson(params)),
		header, hbp.timeout)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))
	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	return &proto.Order{
		OrderID:respmap["data"].(string),
		Amount:util.ToFloat64(amount),
		Price:util.ToFloat64(price),
		Side:side,
		Status:"ok",
	}, nil
}

func (hbp *HuobiPro) Buy(amount, price string, currencyPair string) (*proto.Order, error) {
	return hbp.placeOrder(proto.BUY, amount, price, convertCurrency(currencyPair), "buy-limit")
}

func (hbp *HuobiPro) Sell(amount, price string, currencyPair string) (*proto.Order, error) {
	return hbp.placeOrder(proto.SELL, amount, price, convertCurrency(currencyPair), "sell-limit")
}

func (hbp *HuobiPro) parseOrder(ordmap map[string]interface{}) *proto.Order {
	ord := &proto.Order{
		OrderID:      fmt.Sprintf("%v", ordmap["id"]),
		Amount:       util.ToFloat64(ordmap["amount"]),
		Price:        util.ToFloat64(ordmap["price"]),
		DealedAmount: util.ToFloat64(ordmap["field-amount"]),
		Fee:          util.ToFloat64(ordmap["field-fees"]),
		OrderTime:    fmt.Sprintf("%v", ordmap["created-at"]),
	}

	state := ordmap["state"].(string)
	switch state {
	case "submitted":
		ord.Status = proto.ORDER_UNFINISH
	case "filled":
		ord.Status = proto.ORDER_FINISH
	case "partial-filled":
		ord.Status = proto.ORDER_PART_FINISH
	case "canceled", "partial-canceled":
		ord.Status = proto.ORDER_CANCEL_ING
	default:
		ord.Status = proto.ORDER_UNFINISH
	}

	//if ord.DealedAmount > 0.0 {
	//	ord.AvgPrice = ToFloat64(ordmap["field-cash-amount"]) / ord.DealAmount
	//}

	typeS := ordmap["type"].(string)
	switch typeS {
	case "buy-limit":
		ord.Side = proto.BUY
	case "sell-limit":
		ord.Side = proto.SELL
	}
	return ord
}

func (hbp *HuobiPro) GetOneOrder(orderId string, currency string) (*proto.Order, error) {
	path := "/v1/order/orders/" + orderId
	params := url.Values{}
	hbp.buildPostForm("GET", path, &params)
	resp, err := util.HttpGet(BASE_URL+path+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, respmap)
	if err != nil {
		return nil, err
	}
	if respmap["status"].(string) != "ok" {
		return nil, errors.New(respmap["err-code"].(string))
	}

	datamap := respmap["data"].(map[string]interface{})
	order := hbp.parseOrder(datamap)
	order.Currency = currency
	//log.Println(respmap)
	return &proto.Order{}, nil
}

func (hbp *HuobiPro) CancelOrder(orderId string, currencyPair string) (bool, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	params := url.Values{}
	hbp.buildPostForm("POST", path, &params)
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("Accept-Language", "zh-cn")
	resp, err := util.Request("POST", BASE_URL+path+"?"+params.Encode(), "application/json", strings.NewReader(hbp.toJson(params)),
		header, hbp.timeout)
	if err != nil {
		return false, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	if respmap["status"].(string) != "ok" {
		return false, errors.New(string(resp))
	}

	return true, nil
}

func (hbp *HuobiPro) GetExchangeName() string {
	return "huobi.com"
}

func (hbp *HuobiPro) GetTicker(currencyPair string) (float64, error) {
	_ = BASE_URL + "/market/detail/merged?symbol=" + currencyPair

	return 0.0, nil
}

func (hbp *HuobiPro) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := BASE_URL + "/market/depth?symbol=%s&type=step0"
	resp, err := util.HttpGet(fmt.Sprintf(url, convertCurrency(currencyPair)), nil)
	if err != nil {
		return nil, err
	}

	body := Depth{}
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Huobipro, currencyPair, err)
	}
	if body.Status == "error" {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Huobipro, currencyPair, body.Status)
	}

	asks := body.Tick.Asks
	bids := body.Tick.Bids
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

func priceOfDepth(terms [][]float64, depth float64) (float64, float64, error) {
	var sum float64 = 0.0
	for _, term := range terms {
		p := term[0]
		c := term[1]
		total := p * c
		sum += total
		if sum > float64(depth) {
			return p, sum, nil
		}
	}
	return 0.0, 0.0, fmt.Errorf("has no enough depth sum:%v depth:%v", sum, depth)
}

func (hbp *HuobiPro) Deposit(currency string, amount float64) error {
	return nil
}
func (hbp *HuobiPro) Withdraw(currency string, amount float64) error {
	return nil
}

func (hbp *HuobiPro) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hbp.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(BASE_URL, "https://", "", len(BASE_URL))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := util.SHA256Base64Sign(hbp.secretKey, payload)
	postForm.Set("Signature", sign)
	return nil
}

func (hbp *HuobiPro) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}

func convertCurrency(currency string) string {
	return strings.ToLower(strings.Replace(currency, "_","",-1))
}