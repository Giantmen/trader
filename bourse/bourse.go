package bourse

import "github.com/Giantmen/trader/proto"

type Bourse interface {
	GetTicker(currencyPair string) (float64, error)
	GetPriceOfDepth(size, depth int, currencyPair string) (*proto.Price, error)
	GetAccount() (*proto.Account, error)
	Sell(amount, price string, currencyPair string) (*proto.Order, error)
	Buy(amount, price string, currencyPair string) (*proto.Order, error)
	CancelOrder(orderId string, currencyPair string) (bool, error)
	GetOneOrder(orderId string, currencyPair string) (*proto.Order, error)
}
