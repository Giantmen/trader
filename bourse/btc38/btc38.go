package btc38

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

const (
	MARKET_URL = "http://api.btc38.com/v1/"
	//TICKER_API = "ticker?currency=%s"
	DEPTH_API = "depth.php?c=%s&mk_type=cny"

	// TRADE_URL                 = "https://trade.btc38.com/api/"
	// GET_ACCOUNT_API           = "getAccountInfo"
	// GET_ORDER_API             = "getOrder"
	// GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	// CANCEL_ORDER_API          = "cancelOrder"
	// PLACE_ORDER_API           = "order"
	// WITHDRAW_API              = "withdraw"
	// CANCELWITHDRAW_API        = "cancelWithdraw"
)

type Btc38 struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewBtc38(accessKey, secretKey string, timeout int) (*Btc38, error) {
	return &Btc38{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (btc38 *Btc38) GetTicker(currencyPair string) (float64, error) {
	// url := fmt.Sprintf(MARKET_URL + fmt.Sprintf(TICKER_API, currencyPair))
	// rep, err := util.Request("GET", url, "application/json", nil, nil, btc38.timeout)
	// if err != nil {
	// 	return 0, fmt.Errorf("%s request err %s %v", proto.Btc38, currencyPair, err)
	// }
	// body := Market{}
	// if err := json.Unmarshal(rep, &body); err != nil {
	// 	return 0, fmt.Errorf("%s json Unmarshal err %s %v", proto.Btc38, currencyPair, err)
	// }
	// return strconv.ParseFloat(body.Ticker.Last, 64)
	return 0.0, nil
}

func (btc38 *Btc38) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := MARKET_URL + fmt.Sprintf(DEPTH_API, proto.ETC)
	//fmt.Println("btc38 get depth:", url)
	rep, err := util.Request("GET", url, "application/json", nil, nil, 4)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Btc38, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Btc38, currencyPair, err)
	}

	var sellsum float64
	var sellprice float64
	var buysum float64
	var buyprice float64
	var len int = len(body.Asks)

	for i := len - 1; i >= 0; i-- {
		sellsum += body.Asks[i][1]
		if sellsum > float64(depth) {
			sellprice = body.Asks[i][0]
			break
		}
	}

	for i := 0; i < len; i++ {
		buysum += body.Bids[i][1]
		if buysum > float64(depth) {
			buyprice = body.Bids[i][0]
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

func (btc38 *Btc38) GetAccount() (*proto.Account, error) {
	// params := url.Values{}
	// params.Set("method", "getAccountInfo")
	// btc38.buildPostForm(&params)
	// //log.Println(params.Encode())
	//
	// rep, err := util.Request("POST", TRADE_URL+GET_ACCOUNT_API, "application/x-www-form-urlencoded",
	// 	strings.NewReader(params.Encode()), nil, btc38.timeout)
	// if err != nil {
	// 	return nil, fmt.Errorf("request GetAccount err %v", err)
	// }
	// myaccount := MyAccount{}
	// err = json.Unmarshal(rep, &myaccount)
	// if err != nil {
	// 	return nil, fmt.Errorf("json Unmarshal err %v", err)
	// }
	//
	// account := proto.Account{}
	// account.Asset = myaccount.Result.NetAssets
	// account.Bourse = strings.ToLower(proto.Btc38)
	// account.SubAccounts = make(map[string]proto.SubAccount)
	//
	// //cny
	// subAcc := proto.SubAccount{}
	// subAcc.Available = myaccount.Result.Balance.CNY.Amount
	// subAcc.Forzen = myaccount.Result.Frozen.CNY.Amount
	// subAcc.Currency = strings.ToLower(myaccount.Result.Balance.CNY.Currency)
	// account.SubAccounts[subAcc.Currency] = subAcc
	// //btc
	// subAcc = proto.SubAccount{}
	// subAcc.Available = myaccount.Result.Balance.BTC.Amount
	// subAcc.Forzen = myaccount.Result.Frozen.BTC.Amount
	// subAcc.Currency = strings.ToLower(myaccount.Result.Balance.BTC.Currency)
	// account.SubAccounts[subAcc.Currency] = subAcc
	// //etc
	// subAcc = proto.SubAccount{}
	// subAcc.Available = myaccount.Result.Balance.ETC.Amount
	// subAcc.Forzen = myaccount.Result.Frozen.ETC.Amount
	// subAcc.Currency = strings.ToLower(myaccount.Result.Balance.ETC.Currency)
	// account.SubAccounts[subAcc.Currency] = subAcc
	// //eth
	// subAcc = proto.SubAccount{}
	// subAcc.Available = myaccount.Result.Balance.ETH.Amount
	// subAcc.Forzen = myaccount.Result.Frozen.ETH.Amount
	// subAcc.Currency = strings.ToLower(myaccount.Result.Balance.ETH.Currency)
	// account.SubAccounts[subAcc.Currency] = subAcc
	// //ltc
	// subAcc = proto.SubAccount{}
	// subAcc.Available = myaccount.Result.Balance.LTC.Amount
	// subAcc.Forzen = myaccount.Result.Frozen.LTC.Amount
	// subAcc.Currency = strings.ToLower(myaccount.Result.Balance.LTC.Currency)
	// account.SubAccounts[subAcc.Currency] = subAcc
	//
	// return &account, nil
	return nil, nil
}

func (btc38 *Btc38) placeOrder(side int, amount, price, currencyPair string) (*proto.Order, error) {
	// params := url.Values{}
	// params.Set("method", "order")
	// params.Set("price", price)
	// params.Set("amount", amount)
	// params.Set("currency", currencyPair)
	// params.Set("tradeType", fmt.Sprintf("%d", side))
	// btc38.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+PLACE_ORDER_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, btc38.timeout)
	// var witchside string
	// if side == 1 {
	// 	witchside = proto.BUY
	// } else {
	// 	witchside = proto.SELL
	// }
	// if err != nil {
	// 	return nil, fmt.Errorf("request %s err:%v", witchside, err)
	// }
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
	// return btc38.parseOrder(&myorder)
	//log.Debug("order price:", order.Price, "send price:", price) //对比执行完订单和下发的区别

	return nil, nil
}

func (btc38 *Btc38) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	//return btc38.placeOrder(proto.BUY_N, amount, price, currencyPair)
	return nil, nil
}

func (btc38 *Btc38) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	//return btc38.placeOrder(proto.SELL_N, amount, price, currencyPair)
	return nil, nil
}

func (btc38 *Btc38) CancelOrder(orderId, currencyPair string) (bool, error) {
	// params := url.Values{}
	// params.Set("method", "cancelOrder")
	// params.Set("id", orderId)
	// params.Set("currency", currencyPair)
	// btc38.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+CANCEL_ORDER_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, btc38.timeout)
	// if err != nil {
	// 	return false, fmt.Errorf("request CancelOrder err %v", err)
	// }
	//
	// body := Respons{}
	// err = json.Unmarshal(rep, &body)
	// if err != nil {
	// 	return false, fmt.Errorf("json Unmarshal err %v", err)
	// }
	// if body.Code == 1000 {
	// 	return true, nil
	// }
	// return false, fmt.Errorf("orderid:%s err:%s", orderId, body.Message)

	return false, nil
}

func (btc38 *Btc38) parseOrder(myorder *MyOrder) (*proto.Order, error) {
	// var status string
	// switch myorder.Status {
	// case 0:
	// 	status = proto.ORDER_UNFINISH
	// case 1:
	// 	status = proto.ORDER_CANCEL
	// case 2:
	// 	status = proto.ORDER_FINISH
	// case 3:
	// 	status = proto.ORDER_PART_FINISH
	// }
	// var Side string
	// if myorder.Type == 1 {
	// 	Side = proto.BUY
	// } else {
	// 	Side = proto.SELL
	// }
	// return &proto.Order{
	// 	Amount:       myorder.TradeAmount,
	// 	Fee:          myorder.Fees,
	// 	OrderID:      myorder.ID,
	// 	Price:        float64(myorder.Price),
	// 	DealedAmount: myorder.TradeAmount,
	// 	Currency:     myorder.Currency,
	// 	Status:       status,
	// 	OrderTime:    time.Now().Format(proto.LocalTime),
	// 	Side:         Side,
	// }, nil

	return nil, nil
}

func (btc38 *Btc38) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	// params := url.Values{}
	// params.Set("method", "getOrder")
	// params.Set("id", orderId)
	// params.Set("currency", currencyPair)
	// btc38.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+GET_ORDER_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, btc38.timeout)
	// if err != nil {
	// 	return nil, fmt.Errorf("request GetOneOrder err %v %s", err, orderId)
	// }
	// myorder := MyOrder{}
	// if err := json.Unmarshal(rep, &myorder); err != nil {
	// 	return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	// }
	// return btc38.parseOrder(&myorder)

	return nil, nil
}

func (btc38 *Btc38) GetUnfinishOrders(currencyPair string) (*[]proto.Order, error) {
	// params := url.Values{}
	// params.Set("method", "getUnfinishedOrdersIgnoreTradeType")
	// params.Set("currency", currencyPair)
	// params.Set("pageIndex", "1")
	// params.Set("pageSize", "100")
	// btc38.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+GET_UNFINISHED_ORDERS_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, btc38.timeout)
	// if err != nil {
	// 	return nil, fmt.Errorf("request GetUnfinishOrders err %s", err)
	// }
	//
	// myorders := []MyOrder{}
	// if err := json.Unmarshal(rep, &myorders); err != nil {
	// 	return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	// }
	// orders := []proto.Order{}
	// for _, myorder := range myorders {
	// 	if order, err := btc38.parseOrder(&myorder); err != nil {
	// 		orders = append(orders, *order)
	// 	}
	// }
	return nil, nil
}

func (btc38 *Btc38) buildPostForm(postForm *url.Values) error {
	// postForm.Set("accesskey", btc38.accessKey)
	// payload := postForm.Encode()
	// secretkeySha, _ := util.SHA1(btc38.secretKey)
	//
	// sign, err := util.MD5Sign(secretkeySha, payload)
	// if err != nil {
	// 	return err
	// }
	// postForm.Set("sign", sign)
	// postForm.Set("reqTime", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	return nil
}
