package huobiO

import "github.com/Giantmen/trader/proto"

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
	return nil, nil
}
func (huobi *Huobi) GetAccount() (*proto.Account, error) {
	return nil, nil
}

func (huobi *Huobi) GetTicker(currencyPair string) (float64, error) {
	return 0, nil
}
func (huobi *Huobi) Sell(amount, price, currencyPair string) (*proto.Order, error) {
	return nil, nil
}
func (huobi *Huobi) Buy(amount, price, currencyPair string) (*proto.Order, error) {
	return nil, nil
}
func (huobi *Huobi) CancelOrder(orderId, currencyPair string) (bool, error) {
	return false, nil
}
func (huobi *Huobi) GetOneOrder(orderId, currencyPair string) (*proto.Order, error) {
	return nil, nil
}
