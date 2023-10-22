package main

import (
	"fmt"

	stan "github.com/nats-io/stan.go"
)

var example = []string{`{
		"order_uid": "b563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
		  "name": "Test Testov",
		  "phone": "+9720000000",
		  "zip": "2639809",
		  "city": "Kiryat Mozkin",
		  "address": "Ploshad Mira 15",
		  "region": "Kraiot",
		  "email": "test@gmail.com"
		},
		"payment": {
		  "transaction": "b563feb7b2b84b6test",
		  "request_id": "",
		  "currency": "USD",
		  "provider": "wbpay",
		  "amount": 1817,
		  "payment_dt": 1637907727,
		  "bank": "alpha",
		  "delivery_cost": 1500,
		  "goods_total": 317,
		  "custom_fee": 0
		},
		"items": [
		  {
			"chrt_id": 9934930,
			"track_number": "WBILMTESTTRACK",
			"price": 453,
			"rid": "ab4219087a764ae0btest",
			"name": "Mascaras",
			"sale": 30,
			"size": "0",
			"total_price": 317,
			"nm_id": 2389212,
			"brand": "Vivienne Sabo",
			"status": 202
		  }
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": "9",
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": "1"
	}`,
	`{
		"order_uid": "firsttest1234",
		"track_number": "testTRACK",
		"entry": "test",
		"delivery": {
		  "name": "Sasha Sasha",
		  "phone": "+9720123000",
		  "zip": "2639101",
		  "city": "Angarsk City",
		  "address": "Drobi 123",
		  "region": "Irkutskaya oblast",
		  "email": "test@yandex.com"
		},
		"payment": {
		  "transaction": "a623feb7b2b84b6test",
		  "request_id": "",
		  "currency": "USD",
		  "provider": "wbpay",
		  "amount": 672,
		  "payment_dt": 127907727,
		  "bank": "delta",
		  "delivery_cost": 120,
		  "goods_total": 317,
		  "custom_fee": 0
		},
		"items": [
		  {
			"chrt_id": 123,
			"track_number": "testTRACK",
			"price": 21,
			"rid": "bc4219087a764ae0btest",
			"name": "Vape",
			"sale": 10,
			"size": "0",
			"total_price": 42,
			"nm_id": 4219212,
			"brand": "Rodnoe",
			"status": 202
		  },
		  {
			"chrt_id": 134,
			"track_number": "testTRACK",
			"price": 34,
			"rid": "cd4219087a764ae0btest",
			"name": "Shisha",
			"sale": 15,
			"size": "0",
			"total_price": 317,
			"nm_id": 2319212,
			"brand": "Nashe",
			"status": 202
		  }
		],
		"locale": "ru",
		"internal_signature": "",
		"customer_id": "test1",
		"delivery_service": "sdek",
		"shardkey": "3",
		"sm_id": 101,
		"date_created": "2022-11-26T06:22:19Z",
		"oof_shard": "1"
	}`,
	"Hello, World!",
}

func main() {
	conn, err := stan.Connect("test-cluster", "test", stan.NatsURL("http://localhost:4222"))
	if err != nil {
		panic(fmt.Sprintf("could not connect to NATS Streaming Server: %s", err))
	}
	for _, s := range example {
		err = conn.Publish("test-subject", []byte(s))
		if err != nil {
			panic(fmt.Sprintf("could not publish message: %s", err))
		}
	}
}
