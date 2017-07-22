package huobi

import (
	"fmt"

	"github.com/Giantmen/trader/bourse/huobi/huobiN"
	"github.com/Giantmen/trader/bourse/huobi/huobiO"
	"github.com/Giantmen/trader/proto"
)

type Huobi struct {
	huobiN *huobiN.Huobi
	huobiO *huobiO.Huobi
}

func NewHuobi(accountid, accessKey, secretKey string, timeout int) (*Huobi, error) {
	huobiN, err := huobiN.NewHuobi(accountid, accessKey, secretKey, timeout)
	if err != nil {
		return nil, err
	}
	huobiO, err := huobiO.NewHuobi(accountid, accessKey, secretKey, timeout)
	if err != nil {
		return nil, err
	}
	return &Huobi{
		huobiN: huobiN,
		huobiO: huobiO,
	}, nil
}

func (huobi *Huobi) GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error) {
	switch currencyPair {
	case proto.BTC_CNY, proto.LTC_CNY:
		return huobi.huobiO.GetPriceOfDepth(size, depth, currencyPair)
	case proto.ETH_CNY, proto.ETC_CNY:
		return huobi.huobiN.GetPriceOfDepth(size, depth, currencyPair)
	default:
		return nil, fmt.Errorf("currencyPair err %s", currencyPair)
	}
}

func (huobi *Huobi) GetAccount() (*proto.Account, error) {

	//huobi.huobiO.GetAccount(size, depth, currencyPair)
	return huobi.huobiN.GetAccount()
}

func (huobi *Huobi) GetTicker(currencyPair string) (float64, error) {
	return 0, nil
}
func (huobi *Huobi) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return nil, nil
}

func (huobi *Huobi) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	switch currencyPair {
	case proto.BTC_CNY, proto.LTC_CNY:
		return huobi.huobiO.Buy(amount, price, currencyPair)
	case proto.ETH_CNY, proto.ETC_CNY:
		return huobi.huobiN.Buy(amount, price, currencyPair)
	default:
		return nil, fmt.Errorf("currencyPair err %s", currencyPair)
	}
}

func (huobi *Huobi) CancelOrder(orderId, currencyPair string) (bool, error) {
	return false, nil
}
func (huobi *Huobi) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	return nil, nil
}
