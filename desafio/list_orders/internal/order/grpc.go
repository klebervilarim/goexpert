package order

import (
	"context"

	"orders-app/proto/orderpb"
)

type GRPCServer struct {
	orderpb.UnimplementedOrderServiceServer
	Repo *Repository
}

func (s *GRPCServer) ListOrders(ctx context.Context, _ *orderpb.Empty) (*orderpb.OrderList, error) {
	orders, err := s.Repo.List()
	if err != nil {
		return nil, err
	}

	var pbOrders []*orderpb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &orderpb.Order{
			Id:       int32(o.ID),
			Customer: o.Customer,
			Amount:   o.Amount,
		})
	}

	return &orderpb.OrderList{Orders: pbOrders}, nil
}
