package bourse

import "github.com/Giantmen/trader/proto"

type Bourse interface {
	GetTicker(currencyPair string) (float64, error)
	GetPriceOfDepth(size int, depth float64, currencyPair string) (*proto.Price, error)
	GetAccount() (*proto.Account, error)
	Sell(amount, price, currencyPair string) (*proto.Order, error)
	Buy(amount, price, currencyPair string) (*proto.Order, error)
	CancelOrder(orderId, currencyPair string) (bool, error)
	GetOneOrder(orderId, currencyPair string) (*proto.Order, error)
}
