package bitfinex

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_URL = "https://api.bitfinex.com/v1"
)

type Bitfinex struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewBitfinex(accessKey, secretKey string, timeout int) (*Bitfinex,error) {
	return &Bitfinex{
		accessKey,
		secretKey,
		timeout,
	},nil
}

func (bfx *Bitfinex) GetTicker(currencyPair string) (float64, error) {
	//pubticker
	currencyPair = bfx.convertCurrency(currencyPair)

	apiUrl := fmt.Sprintf("%s/pubticker/%s", BASE_URL, currencyPair)
	resp, err := util.HttpGet(apiUrl, nil)
	if err != nil {
		return 0.0, err
	}

	var bodyDataMap map[string]interface{}
	//fmt.Printf("\n%s\n", respData);
	err = json.Unmarshal(resp, &bodyDataMap)
	if err != nil {
		fmt.Println(string(resp))
		return 0.0, err
	}

	if bodyDataMap["error"] != nil {
		return 0.0, errors.New(bodyDataMap["error"].(string))
	}

	return util.ToFloat64(bodyDataMap["last_price"]), nil
}

func (bfx *Bitfinex) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	apiUrl := fmt.Sprintf("%s/book/%s?limit_bids=%d&limit_asks=%d", BASE_URL, bfx.convertCurrency(currencyPair), size, size)
	fmt.Println(apiUrl)
	resp, err := util.HttpGet(apiUrl, nil)
	if err != nil {
		return nil, err
	}
	println("resp:", resp)

	body := Depth{}
	body.Asks = make([]SubDepth,0)
	body.Bids = make([]SubDepth,0)
	if err := json.Unmarshal(resp, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bitfinex, currencyPair, err)
	}
	bids := body.Bids
	asks := body.Asks

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

func priceOfDepth(terms []SubDepth, depth float64) (float64, float64, error) {
	var sum float64 = 0.0
	for _, term := range terms {
		//entry, _ := term.([]interface{})
		//p := entry[0].(string)
		//c := entry[1].(string)
		price, err := strconv.ParseFloat(term.Price, 64)
		if err != nil {
			return 0.0, 0.0, err
		}
		amount, err := strconv.ParseFloat(term.Amount, 64)
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

//func (bfx *Bitfinex) GetWalletBalances() (map[string]*Account, error) {
//	var respmap []interface{}
//	err := bfx.doAuthenticatedRequest("GET", "balances", map[string]interface{}{}, &respmap)
//	if err != nil {
//		return nil, err
//	}
//	//log.Println(respmap)
//
//	walletmap := make(map[string]*Account, 1)
//
//	for _, v := range respmap {
//		subacc := v.(map[string]interface{})
//		typeStr := subacc["type"].(string)
//
//		currency := NewCurrency(subacc["currency"].(string), "")
//
//		if currency == UNKNOWN {
//			continue
//		}
//
//		//typeS := subacc["type"].(string)
//		amount := ToFloat64(subacc["amount"])
//		available := ToFloat64(subacc["available"])
//
//		account := walletmap[typeStr]
//		if account == nil {
//			account = new(Account)
//			account.SubAccounts = make(map[Currency]SubAccount, 6)
//		}
//
//		account.NetAsset = amount
//		account.Asset = amount
//		account.SubAccounts[currency] = SubAccount{
//			Currency:     currency,
//			Amount:       available,
//			ForzenAmount: amount - available,
//			LoanAmount:   0}
//
//		walletmap[typeStr] = account
//	}
//
//	return walletmap, nil
//}

/*defalut only return exchange wallet balance*/
func (bfx *Bitfinex) GetAccount() (*proto.Account, error) {
	params := make(map[string]interface{})
	nonce := time.Now().UnixNano()
	params["request"]= "/v1/balances"
	params["nonce"]=fmt.Sprintf("%d.2", nonce)

	p, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(p)
	sign, _ := util.Sha384Sign(bfx.secretKey, encoded)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	h.Add("X-BFX-APIKEY", bfx.accessKey)
	h.Add("X-BFX-PAYLOAD", encoded)
	h.Add("X-BFX-SIGNATURE", sign)
	resp, err := util.HttpGet(BASE_URL+"/"+"balances", &h)
	if err != nil {
		return nil, err
	}

	var dataMap []MyAccount
	err = json.Unmarshal(resp, &dataMap)
	if err != nil {
		return nil, err
	}

	account := &proto.Account{}
	account.SubAccounts = make(map[string]proto.SubAccount)

	for _, v := range dataMap {
		currency := v.Currency
		account.SubAccounts[currency] = proto.SubAccount{
			Currency:  currency,
			Available: util.ToFloat64(v.Available),
			//Forzen:    util.ToFloat64(vv["locked"]),
		}
	}
	account.Bourse= proto.Bitfinex

	return account, nil
}

func (bfx *Bitfinex) placeOrder(side, amount, price string, currencyPair string) (*proto.Order, error) {
	path := "order/new"
	nonce := time.Now().UnixNano()

	params := map[string]interface{}{
		"symbol":   bfx.convertCurrency(currencyPair),
		"amount":   amount,
		"price":    price,
		"side":     side,
		"type":     "exchange limit",
		"exchange": "bitfinex",
		"request":  "/v1/"+path,
		"nonce":    fmt.Sprintf("%d.2", nonce),
	}

	p, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(p)
	sign, _ := util.Sha384Sign(bfx.secretKey, encoded)
	h := http.Header{}
	h.Set("Accept", "application/json")
	h.Add("X-BFX-APIKEY", bfx.accessKey)
	h.Add("X-BFX-PAYLOAD", encoded)
	h.Add("X-BFX-SIGNATURE", sign)

	resp, err := util.Request("POST", BASE_URL+"/"+path, "application/json", strings.NewReader(string(p)), h, bfx.timeout)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))
	respmap:=make(map[string]interface{},0)
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return nil, err
	}

	return &proto.Order{
		OrderID: fmt.Sprintf("%d",respmap["id"]),
		Price:        util.ToFloat64(price),
		Amount:       util.ToFloat64(amount),
		DealedAmount: util.ToFloat64(respmap["executed_amount"]),

		Currency: currencyPair,
		Side:     side,
	}, nil
}

// limit buy
func (bfx *Bitfinex) Buy(amount, price string, currencyPair string) (*proto.Order, error) {
	return bfx.placeOrder("buy", amount, price, currencyPair)
}

func (bfx *Bitfinex) Sell(amount, price string, currencyPair string) (*proto.Order, error) {
	return bfx.placeOrder("sell", amount, price, currencyPair)
}

func (bfx *Bitfinex) CancelOrder(orderId string, currencyPair string) (bool, error) {
	path := "order/cancel"

	nonce := time.Now().UnixNano()

	params := map[string]interface{}{
		"order_id": util.ToUint64(orderId),
		"nonce":    fmt.Sprintf("%d.2", nonce),
		"request":"/v1/"+path,
	}

	p, err := json.Marshal(params)
	if err != nil {
		return false, err
	}

	encoded := base64.StdEncoding.EncodeToString(p)
	sign, _ := util.Sha384Sign(bfx.secretKey, encoded)
	h := http.Header{}
	h.Set("Accept", "application/json")
	h.Add("X-BFX-APIKEY", bfx.accessKey)
	h.Add("X-BFX-PAYLOAD", encoded)
	h.Add("X-BFX-SIGNATURE", sign)
	resp, err := util.Request("POST", BASE_URL+"/"+path, "application/json", strings.NewReader(string(p)), h, bfx.timeout)
	if err != nil {
		return false, err
	}
	fmt.Println(string(resp))
	respmap := make(map[string]interface{},0)
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return false, err
	}

	return respmap["is_cancelled"].(bool), nil
}

func (bfx *Bitfinex) Deposit(currency string, amount float64) error {
	return nil
}
func (bfx *Bitfinex) Withdraw(currency string, amount float64) error {
	return nil
}

func (bfx *Bitfinex) GetOneOrder(orderId string, currencyPair string) (*proto.Order, error) {
	return nil, nil
}

func (bfx *Bitfinex) adaptTimestamp(timestamp string) int {
	times := strings.Split(timestamp, ".")
	intTime, _ := strconv.Atoi(times[0])
	return intTime
}
func (bfx *Bitfinex) convertCurrency(pair string) string {
	return strings.ToUpper(strings.Replace(pair, "_","",-1))
}
