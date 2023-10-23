package pgs

import (
	"context"
	"fmt"

	"github.com/PotterVombad/L0/internal/models"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type (
	Database interface {
		SaveOrder(ctx context.Context, order models.Order) error
		GetAllOrders(ctx context.Context) (map[string]models.Order, error)
	}

	PgsDB struct {
		client *pgx.Conn
	}
)

func (db PgsDB) SaveOrder(
	ctx context.Context,
	o models.Order,
) error {
	if err := db.insertInOrder(ctx, o); err != nil {
		return err
	}
	if err := db.insertInPayment(ctx, o); err != nil {
		return err
	}
	if err := db.insertInDelivery(ctx, o); err != nil {
		return err
	}
	if err := db.insertInItem(ctx, o); err != nil {
		return err
	}

	return nil
}

func (db PgsDB) insertInOrder(ctx context.Context, o models.Order) error {
	q := "INSERT INTO orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	raw, err := db.client.Query(
		ctx, q,
		o.Uid, o.TrackNumber, o.Entry, o.Local,
		o.InternalSignature, o.CustomId, o.DeliveryService,
		o.Shardkey, o.SmId, o.DateCreated, o.OofShard)
	raw.Close()
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}
	return nil
}

func (db PgsDB) insertInPayment(ctx context.Context, o models.Order) error {
	q := "INSERT INTO payment VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	raw, err := db.client.Query(ctx, q,
		o.Payment.Transaction, o.Uid, o.Payment.RequestId,
		o.Payment.Currency, o.Payment.Provider, o.Payment.Amount,
		o.Payment.PaymentDt, o.Payment.Bank, o.Payment.DeliveryCost,
		o.Payment.GoodsTotal, o.Payment.CustomFee)
	raw.Close()
	if err != nil {
		return fmt.Errorf("insert in payment: %w", err)
	}
	return nil
}

func (db PgsDB) insertInDelivery(ctx context.Context, o models.Order) error {
	q := "INSERT INTO delivery (Orders_Uid, Name, Phone, Zip, City, Address, Region, Email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	raw, err := db.client.Query(ctx, q,
		o.Uid, o.Delivery.Name, o.Delivery.Phone,
		o.Delivery.Zip, o.Delivery.City, o.Delivery.Address,
		o.Delivery.Region, o.Delivery.Email)
	raw.Close()
	if err != nil {
		return fmt.Errorf("insert in delivery: %w", err)
	}
	return nil
}

func (db PgsDB) insertInItem(ctx context.Context, o models.Order) error {
	q := "INSERT INTO item VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	for _, item := range o.Items {
		raw, err := db.client.Query(ctx, q,
			item.ChrtId, o.Uid, item.TrackNumber,
			item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
		raw.Close()
		if err != nil {
			return fmt.Errorf("insert in item: %w", err)
		}
	}
	return nil
}

func (db PgsDB) GetAllOrders(
	ctx context.Context,
) (map[string]models.Order, error) {
	q := `
	SELECT o.uid, o.track_number, o.entry, o.local, o.internal_signature, 
	o.custom_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard, 
	p.transaction, p.request_id, p.currency, p.provider, p.amount,
	p.payment_date, p.bank, p.delivery_cost, p.goods_total, p.custom_fee, 
	d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,	
	i.chrt_id, i.track_number, i.price, i.r_id, i.name, i.sale, 
	i.size, i.total_price, i.nm_id, i.brand, i.status

		FROM orders as o
		JOIN payment as p
		ON p.orders_uid = o.uid
		JOIN delivery as d
		ON d.orders_uid = o.uid
		JOIN item as i
		ON i.orders_uid = o.uid
		order by o.uid`

	rows, err := db.client.Query(context.Background(), q)
	if err != nil {
		return nil, fmt.Errorf("querying table: %w", err)
	}
	defer rows.Close()

	orders := make(map[string]models.Order)
	var firstOrder models.Order
	for rows.Next() {
		order := models.Order{}
		item := models.Item{}
		err := rows.Scan(
			&order.Uid,
			&order.TrackNumber,
			&order.Entry,
			&order.Local,
			&order.InternalSignature,
			&order.CustomId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard,
			&order.Payment.Transaction,
			&order.Payment.RequestId,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDt,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return make(map[string]models.Order), fmt.Errorf("scanning table rows: %w", err)
		}
		if firstOrder.Uid == "" {
			firstOrder = order
		}
		if firstOrder.Uid != order.Uid {
			orders[firstOrder.Uid] = firstOrder
			firstOrder = order
		}
		firstOrder.Items = append(firstOrder.Items, item)
	}

	orders[firstOrder.Uid] = firstOrder

	return orders, nil
}

func (db PgsDB) Close(ctx context.Context) {
	if err := db.client.Close(ctx); err != nil {
		log.Errorf("didn't close db connection: %v", err)
	}
}

func MustNew(ctx context.Context, url string) PgsDB {
	client, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(fmt.Sprintf("db connection: %s", err))
	}

	if err := client.Ping(ctx); err != nil {
		panic(fmt.Sprintf("ping error: %s", err))
	}

	return PgsDB{
		client: client,
	}
}
