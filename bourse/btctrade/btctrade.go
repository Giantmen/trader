package btctrade

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

var (
	API_URL          = "https://api.btctrade.com/api"
	DEPTH_URL        = "/depth?coin=%s"
	BUY_URL          = "/buy"
	SELL_URL         = "/sell"
	GET_ACCOUNT_API  = "/balance"
	GET_ORDER_API    = "/fetch_order"
	DELETE_ORDER_API = "/cancel_order"
)

type Btctrade struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewBtctrade(accessKey, secretKey string, timeout int) (*Btctrade, error) {
	return &Btctrade{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (b *Btctrade) GetTicker(currencyPair string) (float64, error) {
	return 0, nil
}

// 获取满足某个深度的价格
func (b *Btctrade) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf(API_URL+DEPTH_URL, b.convertCurrencyPair(currencyPair))
	rep, err := util.Request("GET", url, "", nil, nil, b.timeout)
	if err != nil {
		return nil, err
	}

	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Btctrade, currencyPair, err)
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

func (b *Btctrade) GetAccount() (*proto.Account, error) {
	params := url.Values{}
	b.buildPostForm(&params)
	rep, err := util.Request("POST", API_URL+GET_ACCOUNT_API, "application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()), nil, b.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetAccount err %v", err)
	}
	myaccount := MyAccount{}
	err = json.Unmarshal(rep, &myaccount)
	if err != nil {
		return nil, fmt.Errorf("json Unmarshal err %s", err)
	}

	account := proto.Account{}
	account.Asset, _ = strconv.ParseFloat(myaccount.Asset, 64)
	account.Bourse = strings.ToLower(proto.Btctrade)
	account.SubAccounts = make(map[string]proto.SubAccount)

	//cny
	subAcc := proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.CnyBalance, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.CnyReserved, 64)
	subAcc.Currency = proto.CNY
	account.SubAccounts[subAcc.Currency] = subAcc
	//btc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.BtcBalance, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.BtcReserved, 64)
	subAcc.Currency = proto.BTC
	account.SubAccounts[subAcc.Currency] = subAcc
	//etc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.EtcBalance, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.EtcReserved, 64)
	subAcc.Currency = proto.ETC
	account.SubAccounts[subAcc.Currency] = subAcc
	//eth
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.EthBalance, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.EthReserved, 64)
	subAcc.Currency = proto.ETH
	account.SubAccounts[subAcc.Currency] = subAcc
	//ltc
	subAcc = proto.SubAccount{}
	subAcc.Available, _ = strconv.ParseFloat(myaccount.LtcBalance, 64)
	subAcc.Forzen, _ = strconv.ParseFloat(myaccount.LtcReserved, 64)
	subAcc.Currency = proto.LTC
	account.SubAccounts[subAcc.Currency] = subAcc
	return &account, nil
}

func (b *Btctrade) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return b.placeOrder(proto.BUY, amount, price, currencyPair)
}

func (b *Btctrade) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return b.placeOrder(proto.SELL, amount, price, currencyPair)
}

func (b *Btctrade) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("coin", b.convertCurrencyPair(currencyPair))
	params.Set("price", price)
	params.Set("amount", amount)
	b.buildPostForm(&params)
	var SIDE_URL string
	if side == proto.BUY {
		SIDE_URL = BUY_URL
	} else if side == proto.SELL {
		SIDE_URL = SELL_URL
	}
	rep, err := util.Request("POST", API_URL+SIDE_URL,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, b.timeout)
	if err != nil {
		return nil, fmt.Errorf("request %s err:%v", side, err)
	}

	myorder := OrderReply{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if !myorder.Result {
		return nil, fmt.Errorf("%s err %s", side, myorder.Message)
	}

	Price, _ := strconv.ParseFloat(price, 64)
	Amount, _ := strconv.ParseFloat(amount, 64)
	return &proto.Order{
		OrderID:      myorder.ID,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Price:        Price,
		Amount:       Amount,
		DealedAmount: float64(0),
		Fee:          float64(0),
		Status:       proto.ORDER_UNFINISH,
		Currency:     currencyPair,
		Side:         side,
	}, nil
}

func (b *Btctrade) CancelOrder(orderId, currencyPair string) (bool, error) {
	params := url.Values{}
	params.Set("id", orderId)
	b.buildPostForm(&params)
	rep, err := util.Request("POST", API_URL+DELETE_ORDER_API,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, b.timeout)
	if err != nil {
		return false, fmt.Errorf("request CancelOrder err %s", err)
	}
	if rep == nil {
		return false, err
	}
	cancel := OrderReply{}
	if err := json.Unmarshal(rep, &cancel); err != nil {
		return false, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if !cancel.Result {
		return false, fmt.Errorf("CancelOrder err %s", cancel.Message)
	}
	return true, nil
}

func (b *Btctrade) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	params := url.Values{}
	params.Set("id", orderId)
	b.buildPostForm(&params)
	rep, err := util.Request("POST", API_URL+GET_ORDER_API,
		"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
		nil, b.timeout)
	if err != nil {
		return nil, fmt.Errorf("request GetOneOrder err %v", err)
	}

	myorder := MyOrder{}
	if err := json.Unmarshal(rep, &myorder); err != nil {
		return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	}
	if myorder.ID == 0 && !myorder.Result {
		return nil, fmt.Errorf("GetOneOrder err id: %s %s", orderId, myorder.Message)
	}
	order, err := b.parseOrder(&myorder)
	if err != nil {
		return order, err
	}
	order.Currency = currencyPair
	order.Fee = b.convertFee(currencyPair) * order.DealedAmount
	return order, nil
}

func (b *Btctrade) parseOrder(myorder *MyOrder) (*proto.Order, error) {
	var Status string
	switch myorder.Status {
	case "open":
		Status = proto.ORDER_UNFINISH
	case "closed":
		Status = proto.ORDER_FINISH
	case "cancelled":
		Status = proto.ORDER_CANCEL
	}
	return &proto.Order{
		OrderID:      fmt.Sprintf("%d", myorder.ID),
		OrderTime:    myorder.Datetime,
		Price:        myorder.Price,
		Amount:       myorder.AmountOriginal,
		DealedAmount: myorder.Trades.SumNumber,
		Status:       Status,
		Side:         myorder.Type,
	}, nil
}

func (b *Btctrade) buildPostForm(postForm *url.Values) error {
	postForm.Set("key", b.accessKey)
	postForm.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	postForm.Set("version", fmt.Sprintf("%d", 2))
	payload := postForm.Encode()

	md5sign, err := util.MD5(b.secretKey)
	if err != nil {
		return err
	}
	sign, err := util.SHA256Sign(md5sign, payload)
	if err != nil {
		return err
	}
	postForm.Set("signature", sign)
	return nil
}

func (b *Btctrade) convertCurrencyPair(currencyPair string) string {
	switch currencyPair {
	case proto.BTC_CNY:
		return "btc"
	case proto.ETH_CNY:
		return "eth"
	case proto.ETC_CNY:
		return "etc"
	case proto.LTC_CNY:
		return "ltc"
	default:
		return ""
	}
}

func (b *Btctrade) convertFee(currencyPair string) float64 {
	switch currencyPair {
	case proto.BTC_CNY:
		return proto.FEE_Btctrade_btc
	case proto.ETH_CNY:
		return proto.FEE_Btctrade_eth
	case proto.ETC_CNY:
		return proto.FEE_Btctrade_etc
	case proto.LTC_CNY:
		return proto.FEE_Btctrade_ltc
	}
	return 0
}
