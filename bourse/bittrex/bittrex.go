package bittrex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Giantmen/trader/proto"
	"github.com/Giantmen/trader/util"
)

const (
	PUBLICAPI = "https://bittrex.com/api/v1.1/"
)

const (
	BUY              = "market/buylimit"
	SELL             = "market/selllimit"
	ORDERBOOK        = "public/getorderbook"
	ORDERTRADES      = "returnOrderTrades"
	OPENORDERS       = "returnOpenOrders"
	CANCLEORDER      = "cancelOrder"
	COMPLETEBALANCES = "returnAvailableAccountBalances"
)

type Bittrex struct {
	accessKey string
	secretKey string
	timeout   int
}

func NewPoloniex(accessKey, secretKey string, timeout int) (*Bittrex, error) {
	return &Bittrex{
		accessKey: accessKey,
		secretKey: secretKey,
		timeout:   timeout,
	}, nil
}

func (b *Bittrex) GetTicker(currencyPair string) (float64, error) {
	return 0.0, nil
}

// 获取满足某个深度的价格
func (b *Bittrex) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	url := fmt.Sprintf("%s%s?type=both&market=%s&depth=%d", PUBLICAPI, ORDERBOOK, b.convertCurrencyPair(currencyPair), size)
	rep, err := util.Request("GET", url, "", nil, nil, b.timeout)
	if err != nil {
		return nil, fmt.Errorf("%s request err %s %v", proto.Bittrex, currencyPair, err)
	}
	response := Response{}
	if err := json.Unmarshal(rep, &response); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bittrex, currencyPair, err)
	}
	if !response.Success {
		return nil, fmt.Errorf("err: %s", response.Message)
	}
	body := Depth{}
	if err := json.Unmarshal(response.Result, &body); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bittrex, currencyPair, err)
	}

	var sellsum float64
	var sellprice float64
	var buysum float64
	var buyprice float64
	var len int = len(body.Buy)

	for i := 0; i < len; i++ {
		sellsum += body.Sell[i].Quantity // body.Tick.Asks[i][1]
		if sellsum > float64(depth) {
			sellprice = body.Sell[i].Rate
			break
		}
	}
	for i := 0; i < len; i++ {
		buysum += body.Buy[i].Quantity
		if buysum > float64(depth) {
			buyprice = body.Buy[i].Rate
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

func (b *Bittrex) placeOrder(side, amount, price, currencyPair string) (*proto.Order, error) {
	sign, err := b.buildPostForm()
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json;charset=utf-8")
	header.Set("apisign", sign)

	var witchside string
	if side == proto.BUY {
		witchside = BUY
	} else if side == proto.SELL {
		witchside = SELL
	}
	urls := fmt.Sprintf("https://%s%s?market=%s&quantity=%s&rate=%s",
		PUBLICAPI, witchside, b.convertCurrencyPair(currencyPair), amount, price)
	fmt.Println("urls:", urls)

	rep, err := util.Request("GET", urls, "application/json", nil, header, b.timeout)
	if err != nil {
		return nil, fmt.Errorf("request %v err:%v", side, err)
	}

	response := Response{}
	if err = json.Unmarshal(rep, &response); err != nil {
		return nil, fmt.Errorf("%s json Unmarshal err %s %v", proto.Bittrex, currencyPair, err)
	}
	if !response.Success {
		return nil, fmt.Errorf("err: %s", response.Message)
	}

	tresp := Uuid{}
	err = json.Unmarshal(rep, &tresp)
	if err != nil {
		return nil, err
	}
	// if tresp.Error != "" {
	// 	return nil, fmt.Errorf("%s err:%s", side, tresp.Error)
	// }
	return &proto.Order{
		OrderID:      tresp.OrderNumber,
		OrderTime:    time.Now().Format(proto.LocalTime),
		Price:        float64(0),
		Amount:       float64(0),
		DealedAmount: float64(0),
		Fee:          float64(0),
		Status:       proto.ORDER_UNFINISH,
		Currency:     currencyPair,
		Side:         side,
	}, nil
}

func (b *Bittrex) buildPostForm() (string, error) {
	v := url.Values{}
	v.Set("apikey", b.accessKey)
	v.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	payload := v.Encode()
	sign, err := util.SHA512Sign(b.secretKey, payload)
	if err != nil {
		return "", err
	}
	return sign, nil
}

func (b *Bittrex) convertCurrencyPair(currencyPair string) string {
	switch currencyPair {
	case proto.BTC_BTM:
		return "BTC-BTM"
	case proto.BTC_BTS:
		return "BTC-BTS"
	case proto.BTC_CVC:
		return "BTC-CVC"
	case proto.BTC_EOS:
		return "BTC-EOS"
	case proto.BTC_ETC:
		return "BTC-ETC"
	case proto.BTC_ETH:
		return "BTC-ETH"
	case proto.BTC_LTC:
		return "BTC-LTC"
	case proto.BTC_OMG:
		return "BTC-OMG"
	case proto.BTC_PAY:
		return "BTC-PAY"
	case proto.BTC_SC:
		return "BTC-SC"
	case proto.BTC_SNT:
		return "BTC-SNT"
	default:
		return ""
	}
}
