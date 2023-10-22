package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

var tables = []string{
	`CREATE TABLE IF NOT EXISTS orders 
	(
		uid           	   VARCHAR(30) PRIMARY KEY,
		track_number       VARCHAR(30),
		entry              VARCHAR(30),
		local              VARCHAR(30),
		internal_signature VARCHAR(30),
		custom_id          VARCHAR(30),
		delivery_service   VARCHAR(30),
		shardkey           VARCHAR(30),
		sm_id              SERIAL,
		date_created       TIMESTAMP,
		oof_shard          VARCHAR(30)
	);`,
	`CREATE TABLE IF NOT EXISTS payment 
	(
		transaction   VARCHAR(30) Primary KEY,
		orders_uid    VARCHAR(30) REFERENCES "orders" (uid),
		request_id    VARCHAR(30),
		currency      VARCHAR(30),
		provider      VARCHAR(30),
		amount        real,
		payment_date   serial,
		bank          VARCHAR(30),
		delivery_cost real,
		goods_total   INTEGER,
		custom_fee    real
	);`,
	`CREATE TABLE IF NOT EXISTS item 
	(
		chrt_id      INTEGER PRIMARY KEY,
		orders_uid   VARCHAR(30) REFERENCES "orders" (uid),
		track_number VARCHAR(30),
		price        INTEGER,
		r_id          VARCHAR(30),
		name         VARCHAR(30),
		sale         INTEGER,
		size         VARCHAR(30),
		total_price  INTEGER,
		nm_id        INTEGER,
		brand        VARCHAR(30),
		status       INTEGER
	);`,
	`CREATE TABLE IF NOT EXISTS delivery 
	(
		id SERIAL PRIMARY KEY,
		orders_uid   VARCHAR(30) REFERENCES "orders" (uid),
		name   VARCHAR(30),
		phone  VARCHAR(30),
		zip    VARCHAR(30),
		city   VARCHAR(30),
		address VARCHAR(30),
		region VARCHAR(30),
		email  VARCHAR(30)
	);`,
}

func main() {
	ctx := context.Background()
	url := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s",
		"postgres",
		"postgres",
		"postgresL0",
	)

	client, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}

	res := client.Ping(ctx)
	if res != nil {
		panic(res.Error())
	}

	if err = makeTables(ctx, client); err != nil {
		panic(err)
	}

	fmt.Println("good")
}

func makeTables(ctx context.Context, client *pgx.Conn) (err error) {
	batch := &pgx.Batch{}
	for _, table := range tables {
		batch.Queue(table)
	}
	batchResult := client.SendBatch(ctx, batch)
	defer batchResult.Close()
	return nil
}
