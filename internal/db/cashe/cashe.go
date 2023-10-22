package cashe

import (
	"context"
	"errors"
	"fmt"

	"github.com/PotterVombad/L0/internal/models"
)

type (
	Cashe interface {
		GetOrder(ctx context.Context, uid string) (models.Order, error)
		Save(ctx context.Context, order models.Order) error
		Fill(ctx context.Context, cashe map[string]models.Order)
	}

	InMemory struct {
		cashe map[string]models.Order
	}
)

func (c *InMemory) GetOrder(
	ctx context.Context,
	uid string,
) (models.Order, error) {
	o, ok := c.cashe[uid]
	if !ok {
		return o, fmt.Errorf("order not found  | %s", o.Uid)
	}

	return o, nil
}

var ErrIsExist = errors.New("already exist")

func (c *InMemory) Save(ctx context.Context, order models.Order) error {
	if _, ok := c.cashe[order.Uid]; ok {
		return ErrIsExist
	}

	c.cashe[order.Uid] = order

	return nil
}

func (c *InMemory) Fill(
	ctx context.Context,
	data map[string]models.Order,
) {
	c.cashe = data
}

func New() *InMemory {
	return &InMemory{
		cashe: make(map[string]models.Order),
	}
}
