package models

import "time"

type Order struct {
	OrderUID    string    `json:"order_uid" pg:"order_uid,pk"`
	TrackNumber string    `json:"track_number"`
	Entry       string    `json:"entry"`
	DateCreated time.Time `json:"date_created" pg:"date_created"`
	OofShard    string    `json:"oof_shard"`

	Delivery *Delivery `json:"delivery" pg:"-"`
	Payment  *Payment  `json:"payment" pg:"-"`
	Items    []Item    `json:"items" pg:"-"`
}

type Delivery struct {
	ID       int    `pg:"id,pk"`
	OrderUID string `pg:"order_uid" json:"-"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	ID           int       `pg:"id,pk"`
	OrderUID     string    `pg:"order_uid" json:"-"`
	Transaction  string    `json:"transaction"`
	RequestId    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDt    time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

type Item struct {
	ID        int    `pg:"id,pk"`
	OrderUID  string `pg:"order_uid" json:"-"`
	ChrtID    int64  `json:"chrt_id"`
	ProductID int64  `json:"product_id"`
	Price     int    `json:"price"`
	RID       string `json:"rid"`
	Name      string `json:"name"`
	Sale      int    `json:"sale"`
	Size      string `json:"size"`
	Total     int    `json:"total"`
	NmID      int64  `json:"nm_id"`
	Brand     string `json:"brand"`
	Status    int    `json:"status"`
}
