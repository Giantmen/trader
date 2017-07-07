package api

import (
	"fmt"
	"strings"

	"github.com/Giantmen/trader/bourse"
	"github.com/Giantmen/trader/config"
	"github.com/Giantmen/trader/log"
	"github.com/Giantmen/trader/proto"

	"github.com/solomoner/gozilla"
)

func Register(cfg *config.Config) {
	qs, err := NewQueryService(cfg)
	if err != nil {
		log.Error("register err", "NewQueryService", "err", err)
	}
	log.Debug("register", "qs", qs)
	gozilla.RegisterService(qs, "trader")
}

type Query struct {
	Bourses map[string]bourse.Bourse
}

func NewQueryService(cfg *config.Config) (*Query, error) {
	service, err := bourse.NewService(cfg)
	if err != nil {
		log.Error("new bourse service err", err)
		return nil, err
	}
	return &Query{
		Bourses: service.Bourses,
	}, nil
}

func (q *Query) GetPriceOfDepth(ctx *gozilla.Context, r *proto.DepthQuery) (*proto.Price, error) {
	bou, ok := q.Bourses[strings.ToUpper(r.Bourse)]
	if !ok {
		log.Errorf("get %s err", r.Bourse)
		return nil, fmt.Errorf("get %s err", r.Bourse)
	}

	return bou.GetPriceOfDepth(r.Size, r.Depth, proto.ConvertCurrencyPair(r.Currency))
}

func (q *Query) GetAccount(ctx *gozilla.Context, r *proto.AmountQuery) (*proto.AmountReply, error) {
	bou, ok := q.Bourses[strings.ToUpper(r.Bourse)]
	if !ok {
		log.Errorf("get %s err", r.Bourse)
		return nil, fmt.Errorf("get %s err", r.Bourse)
	}
	account, err := bou.GetAccount()
	if err != nil {
		log.Error("GetAccount err", err)
		return nil, err
	}
	var Accounts = make(map[string]proto.SubAccount)
	for _, currency := range r.Accounts {
		if sub, ok := account.SubAccounts[currency]; ok {
			Accounts[currency] = sub
		} else {
			log.Error("can not find", currency)
		}
	}
	return &proto.AmountReply{
		Bourse:   account.Bourse,
		Asset:    account.Asset,
		Accounts: Accounts,
	}, nil
}

func (q *Query) Sell(ctx *gozilla.Context, r *proto.OrderQuery) (*proto.Order, error) {
	bou, ok := q.Bourses[strings.ToUpper(r.Bourse)]
	if !ok {
		log.Errorf("get %s err", r.Bourse)
		return nil, fmt.Errorf("get %s err", r.Bourse)
	}
	order, err := bou.Sell(r.Amount, r.Price, proto.ConvertCurrencyPair(r.Currency))
	if err != nil {
		log.Error("sell err", err)
	}
	return bou.GetOneOrder(order.OrderID, order.Currency)
}

func (q *Query) Buy(ctx *gozilla.Context, r *proto.OrderQuery) (*proto.Order, error) {
	bou, ok := q.Bourses[strings.ToUpper(r.Bourse)]
	if !ok {
		log.Errorf("get %s err", r.Bourse)
		return nil, fmt.Errorf("get %s err", r.Bourse)
	}
	order, err := bou.Buy(r.Amount, r.Price, proto.ConvertCurrencyPair(r.Currency))
	if err != nil {
		log.Error("sell err", err)
	}
	return bou.GetOneOrder(order.OrderID, order.Currency)
}

func (q *Query) CancelOrder(ctx *gozilla.Context, r *proto.OneOrderQuery) (bool, error) {
	bou, ok := q.Bourses[strings.ToUpper(r.Bourse)]
	if !ok {
		log.Errorf("get %s err", r.Bourse)
		return false, fmt.Errorf("get %s err", r.Bourse)
	}
	return bou.CancelOrder(r.OrderID, proto.ConvertCurrencyPair(r.Currency))
}

func (q *Query) GetOneOrder(ctx *gozilla.Context, r *proto.OneOrderQuery) (*proto.Order, error) {
	bou, ok := q.Bourses[strings.ToUpper(r.Bourse)]
	if !ok {
		log.Errorf("get %s err", r.Bourse)
		return nil, fmt.Errorf("get %s err", r.Bourse)
	}
	return bou.GetOneOrder(r.OrderID, proto.ConvertCurrencyPair(r.Currency))
}
