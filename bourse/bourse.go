package bourse

import (
	"strings"

	"github.com/Giantmen/trader/bourse/chbtc"
	"github.com/Giantmen/trader/bourse/yunbi"
	"github.com/Giantmen/trader/config"
	"github.com/Giantmen/trader/proto"
)

type Bourse interface {
	GetTicker(currencyPair string) (float64, error)
	GetPriceOfDepth(size, depth int, currencyPair string) (*proto.Price, error)
	GetAccount() (*proto.Account, error)
	Sell(amount, price string, currencyPair string) (*proto.Order, error)
	Buy(amount, price string, currencyPair string) (*proto.Order, error)
	CancelOrder(orderId string, currencyPair string) (bool, error)
	GetOneOrder(orderId string, currencyPair string) (*proto.Order, error)
}

type Service struct {
	Bourses map[string]Bourse
}

func NewService(cfg *config.Config) (*Service, error) {
	var bourse = make(map[string]Bourse)

	if yunbi, err := yunbi.NewYunBi(&cfg.Yunbi); err != nil {
		return nil, err
	} else {
		bourse[strings.ToUpper(proto.Yunbi)] = yunbi
	}

	if chbtc, err := chbtc.NewChbtc(&cfg.Chbtc); err != nil {
		return nil, err
	} else {
		bourse[strings.ToUpper(proto.Chbtc)] = chbtc
	}
	return &Service{
		Bourses: bourse,
	}, nil
}
