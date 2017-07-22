package huobiN

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

var (
	API_URL          = "https://be.huobi.com"
	TICKER_URL       = "/market/kline?symbol=%s&period=1min"
	DEPTH_URL        = "/market/depth?symbol=%s&type=step1"
	GET_ACCOUNT_API  = "/v1/account/accounts/%s/balance"
	CREATE_ORDER_API = "/v1/order/orders"
	PLACE_ORDER_API  = "/v1/order/orders/%d/place"
	//DELETE_ORDER_API = "order/delete.json"
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

func (huobi *Huobi) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := API_URL + fmt.Sprintf(DEPTH_URL, huobi.convertCurrencyPair(currencyPair))
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
			Sell: sellprice,
			Buy:  buyprice,
		}, nil
	}
	return nil, fmt.Errorf("sum not enough %v %v", sellsum, depth)
}

func (huobi *Huobi) buildPostForm(postForm *url.Values, path string) error {
	postForm.Set("AccessKeyId", huobi.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().Format(proto.UTCTime))
	signBody := fmt.Sprintf("%s%s", path, postForm.Encode())
	fmt.Println("signBody", signBody)

	md5sign, err := util.MD5(huobi.secretKey)
	if err != nil {
		return err
	}
	sign, err := util.SHA256Sign(md5sign, signBody)
	if err != nil {
		return err
	}
	fmt.Println("sign", sign)
	sign2 := base64.StdEncoding.EncodeToString([]byte(sign))
	//sign2 := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(sign)))
	postForm.Set("Signature", sign2)
	return nil
}

// https://be.huobi.com/v1/order/orders?					  AccessKeyId=e2xxxxxx-99xxxxxx-84xxxxxx-7xxxx&order-id=1234567890&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=2017-05-11T15%3A19%3A30&Signature=4F65x5A2bLyMWVQj3Aqp%2BB4w%2BivaA7n5Oi2SuYtCJ9o%3D
// https://be.huobi.com/v1/account/accounts/11765513/balance?AccessKeyId=be8dd339-031f5028-3ce34334-9f54b&Signature=7a5452443274684941776c74427a436a636732364e574d72436142785a2b6e5773373939703455796375343d&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=2017-07-20T18%3A48%3A11&account-id=11765513
func (huobi *Huobi) GetAccount() (*proto.Account, error) {
	params := url.Values{}
	//params.Set("account-id", huobi.accountid)
	path := fmt.Sprintf(GET_ACCOUNT_API, huobi.accountid)
	signurl := fmt.Sprintf("GET\n%s\n%s\n", "be.huobi.com", path)
	huobi.buildPostForm(&params, signurl)
	//log.Println(params.Encode())
	//fmt.Println("test:", urls+params.Encode())

	fmt.Println("test:", API_URL+path+"?"+params.Encode())
	rep, err := util.Request("GET", API_URL+path+"?"+params.Encode(), "application/json",
		nil, nil, huobi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetAccount err %v", err)
	}
	fmt.Println("rep:", string(rep), err)
	myaccount := MyAccount{}
	if err = json.Unmarshal(rep, &myaccount); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v", err)
	}
	if myaccount.Status != "ok" {
		return nil, fmt.Errorf("%s json Unmarshal err %s %s", proto.HuobiN, myaccount.Err_code, myaccount.Err_msg)
	}

	account := proto.Account{}
	// //account.Asset =
	// account.Bourse = strings.ToLower(proto.Huobi)
	// account.SubAccounts = make(map[string]proto.SubAccount)
	// for _, a := range myaccount.Data.List {
	// 	subAcc := proto.SubAccount{}
	// 	if _, ok := account.SubAccounts[a.Currency]; ok {
	// 		subAcc = account.SubAccounts[a.Currency]
	// 	}
	// 	if a.Type == "trade" {
	// 		subAcc.Available, _ = strconv.ParseFloat(a.Balance, 64)
	// 	} else if a.Type == "frozen" {
	// 		subAcc.Forzen, _ = strconv.ParseFloat(a.Balance, 64)
	// 	}
	// 	subAcc.Currency = a.Currency
	// 	account.SubAccounts[a.Currency] = subAcc
	// }
	return &account, nil
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

func (huobi *Huobi) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("account-id", huobi.accountid)
	params.Set("amount", amount)
	params.Set("price", price)
	params.Set("source", "api")
	params.Set("symbol", huobi.convertCurrencyPair(currencyPair))
	params.Set("type", side)
	signurl := fmt.Sprintf("POST\n%s\n%s\n", "be.huobi.com", CREATE_ORDER_API)
	huobi.buildPostForm(&params, signurl)

	rep, err := util.Request("POST", API_URL+CREATE_ORDER_API, "application/json",
		strings.NewReader(params.Encode()), nil, huobi.timeout)
	if err != nil {
		return nil, fmt.Errorf("request %s err:%v", side, err)
	}
	fmt.Println("rep:", string(rep))
	// placeOrder := PlaceOrder{}
	// if err := json.Unmarshal(rep, &placeOrder); err != nil {
	// 	return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	// }
	// if placeOrder.Code != 1000 {
	// 	return nil, fmt.Errorf("%s err:%v", witchside, placeOrder.Message)
	// }
	// myorder := MyOrder{
	// 	ID:       placeOrder.ID,
	// 	Currency: currencyPair,
	// }
	// return chbtc.parseOrder(&myorder)
	return nil, fmt.Errorf("err:%s", "test")
	//log.Debug("order price:", order.Price, "send price:", price) //对比执行完订单和下发的区别
}

func (huobi *Huobi) GetTicker(currencyPair string) (float64, error) {
	return 0, nil
}
func (huobi *Huobi) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return huobi.placeOrder("sell-limit", amount, price, currencyPair)
}
func (huobi *Huobi) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return huobi.placeOrder("buy-limit", amount, price, currencyPair)
}
func (huobi *Huobi) CancelOrder(orderId, currencyPair string) (bool, error) {
	return false, nil
}
func (huobi *Huobi) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	return nil, nil
}
