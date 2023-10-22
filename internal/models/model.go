package models

import (
	"time"
)

type (
	Order struct {
		Uid               string `json:"order_uid"`
		TrackNumber       string `json:"track_number"`
		Entry             string `json:"entry"`
		Delivery          Delivery
		Payment           Payment
		Items             []Item
		Local             string    `json:"locale"`
		InternalSignature string    `json:"internal_signature"`
		CustomId          string    `json:"customer_id"`
		DeliveryService   string    `json:"delivery_service"`
		Shardkey          string    `json:"shardkey"`
		SmId              int       `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OofShard          string    `json:"oof_shard"`
	}

	Delivery struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	}

	Payment struct {
		Transaction  string  `json:"transaction"`
		RequestId    string  `json:"request_id"`
		Currency     string  `json:"currency"`
		Provider     string  `json:"provider"`
		Amount       float32 `json:"amount"`
		PaymentDt    int     `json:"payment_dt"`
		Bank         string  `json:"bank"`
		DeliveryCost float32 `json:"delivery_cost"`
		GoodsTotal   int     `json:"goods_total"`
		CustomFee    float32 `json:"custom_fee"`
	}

	Item struct {
		ChrtId      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		Rid         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmId        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	}
)

// TODO: do better
func (o *Order) DoBetter() []string {
	return nil
}
