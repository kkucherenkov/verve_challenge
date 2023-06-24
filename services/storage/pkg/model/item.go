package model

import (
	"encoding/json"
	"math"
	"time"
)

type Item struct {
	Id             string    `json:"id"`
	Price          float32   `json:"price"`
	ExpirationDate time.Time `json:"expiration_date"`
}

func (item Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id             string  `json:"id"`
		Price          float32 `json:"price"`
		ExpirationDate string  `json:"expiration_date"`
	}{
		Id:             item.Id,
		Price:          float32(math.Ceil(float64(item.Price*100)) / 100),
		ExpirationDate: item.ExpirationDate.Format(time.DateTime),
	})
}
