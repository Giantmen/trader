package bter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

const (
	accessKey = "CC900CA5-C3D5-4ADA-94DD-894A20AAC933"
	secret    = "9f7eb465126524462a24109c3a6b06a107124bf147f16b72c08b906b47a20413"
	BASE_URL  = "https://api.bter.com/api2/1/"
	DEPTH_API = "orderBook/"
)

// Bter bter
type Bter struct {
	accessKey string
	secretKey string
	timeout   int
}

// NewBter init a bter object
func NewBter(accessKey, secretKey string, timeout int) (*Bter, error) {
	return &Bter{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (bter *Bter) GetTicker(currencyPair string) (float64, error) {
	return 0.0, nil
}

// GetPriceOfDepth get the price of x depth
func (bter *Bter) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf("%s%s%s", BASE_URL, DEPTH_API, currencyPair)
	rep, err := util.Request("GET", url, "application/json", nil, nil, bter.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Bter, currencyPair, err)
	}
	body := Depth{}
	if err := json.Unmarshal(rep, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bter, currencyPair, err)
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

// GetAccount 获取用户账户信息
func (bter *Bter) GetAccount() (*proto.Account, error) {
	header := http.Header{}
	header.Add("Key", bter.accessKey)
	sign, err := util.SHA512Sign(bter.secretKey, "")
	if err != nil {
		return nil, err
	}
	header.Add("Sign", sign)

	url := fmt.Sprintf("%s%s", BASE_URL, "private/balances")
	resp, err := util.Request("post", url, "application/x-www-form-urlencoded", nil, header, bter.timeout)
	if err != nil {
		return nil, err
	}

	myaccount := new(MyAccount)
	err = json.Unmarshal(resp, myaccount)
	if myaccount.Result != "true" {
		println(err)
		return nil, err
	}

	account := new(proto.Account)
	account.SubAccounts = make(map[string]proto.SubAccount)
	for k, v := range myaccount.Available {
		sc := new(proto.SubAccount)
		sc.Currency = strings.ToLower(k)
		ava, _ := strconv.ParseFloat(v, 32)
		sc.Available = ava
		account.SubAccounts[strings.ToLower(k)] = *sc
	}

	return account, nil
}

func (bter *Bter) placeOrder(side int, amount, price, currencyPair string) (*proto.Order, error) {
	// params := url.Values{}
	// params.Set("method", "order")
	// params.Set("price", price)
	// params.Set("amount", amount)
	// params.Set("currency", currencyPair)
	// params.Set("tradeType", fmt.Sprintf("%d", side))
	// bter.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+PLACE_ORDER_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, bter.timeout)
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
	// return bter.parseOrder(&myorder)
	//log.Debug("order price:", order.Price, "send price:", price) //对比执行完订单和下发的区别

	return nil, nil
}

// Buy do buy
func (bter *Bter) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	//return bter.placeOrder(proto.BUY_N, amount, price, currencyPair)
	return nil, nil
}

// Sell do sell
func (bter *Bter) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	//return bter.placeOrder(proto.SELL_N, amount, price, currencyPair)
	return nil, nil
}

// CancelOrder cancle an order
func (bter *Bter) CancelOrder(orderID, currencyPair string) (bool, error) {
	params := url.Values{}
	params.Set("orderNumber", orderID)
	params.Set("currencyPair", currencyPair)

	payload := params.Encode()
	sign, _ := util.SHA512Sign(bter.secretKey, payload)

	header := http.Header{}
	header.Add("Key", bter.accessKey)
	header.Add("Sign", sign)

	url := fmt.Sprintf("%s%s", BASE_URL, "private/cancelOrder")
	body := strings.NewReader(payload)
	resp, err := util.Request("post", url, "application/x-www-form-urlencoded", body, header, bter.timeout)
	if err != nil {
		return false, err
	}

	cancleResp := new(CancleResponse)
	err = json.Unmarshal(resp, cancleResp)
	if err != nil {
		return false, err
	}

	if cancleResp.Result != "true" {
		return false, errors.New(cancleResp.Message)
	}
	return true, nil
}

func (bter *Bter) parseOrder(myorder *MyOrder) (*proto.Order, error) {
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

// GetOneOrder get one order
func (bter *Bter) GetOneOrder(orderID, currencyPair string) (*proto.Order, error) {
	// params := url.Values{}
	// params.Set("method", "getOrder")
	// params.Set("id", orderId)
	// params.Set("currency", currencyPair)
	// bter.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+GET_ORDER_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, bter.timeout)
	// if err != nil {
	// 	return nil, fmt.Errorf("request GetOneOrder err %v %s", err, orderId)
	// }
	// myorder := MyOrder{}
	// if err := json.Unmarshal(rep, &myorder); err != nil {
	// 	return nil, fmt.Errorf("json Unmarshal err %v %s", err, string(rep))
	// }
	// return bter.parseOrder(&myorder)

	return nil, nil
}

// GetUnfinishOrders get unfinished orders
func (bter *Bter) GetUnfinishOrders(currencyPair string) (*[]proto.Order, error) {
	// params := url.Values{}
	// params.Set("method", "getUnfinishedOrdersIgnoreTradeType")
	// params.Set("currency", currencyPair)
	// params.Set("pageIndex", "1")
	// params.Set("pageSize", "100")
	// bter.buildPostForm(&params)
	//
	// rep, err := util.Request("POST", TRADE_URL+GET_UNFINISHED_ORDERS_API,
	// 	"application/x-www-form-urlencoded", strings.NewReader(params.Encode()),
	// 	nil, bter.timeout)
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
	// 	if order, err := bter.parseOrder(&myorder); err != nil {
	// 		orders = append(orders, *order)
	// 	}
	// }
	return nil, nil
}

// func (bter *Bter) buildPostForm(postForm *url.Values) error {
// 	payload := postForm.Encode()
// 	secretkeySha, _ := util.SHA512Sign(bter.secretKey,payload)
//
// 	postForm.Set("sign", sign)
// 	return nil
// }
