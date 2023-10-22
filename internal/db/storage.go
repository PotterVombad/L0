package db

import (
	"context"
	"fmt"

	"github.com/PotterVombad/L0/internal/db/cache"
	"github.com/PotterVombad/L0/internal/db/pgs"
	"github.com/PotterVombad/L0/internal/models"
)

type (
	Store interface {
		Save(ctx context.Context, order models.Order) error
		Get(ctx context.Context, uid string) (order models.Order, err error)
		FillCashe(ctx context.Context) error
	}

	Storage struct {
		cashe cache.Cache
		db    pgs.Database
	}
)

func (s *Storage) Save(
	ctx context.Context,
	order models.Order,
) error {
	if err := s.cashe.Save(ctx, order); err != nil {
		return fmt.Errorf("save to cache | %s", err)
	}

	if err := s.db.SaveOrder(ctx, order); err != nil {
		return fmt.Errorf("save to db  | %s", err)
	}

	return nil
}

func (s *Storage) FillCashe(ctx context.Context) error {
	cashe, err := s.db.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("cashe is empty: %v", err)
	}

	s.cashe.Fill(ctx, cashe)

	return nil
}

func (s *Storage) Get(ctx context.Context, uid string) (order models.Order, err error) {
	return s.cashe.GetOrder(ctx, uid)
}

func MustNew(
	ctx context.Context,
	cashe cache.Cache,
	db pgs.Database,
) Storage {

	s := Storage{
		cashe: cashe,
		db:    db,
	}

	if err := s.FillCashe(ctx); err != nil {
		panic(err)
	}

	return s
}
