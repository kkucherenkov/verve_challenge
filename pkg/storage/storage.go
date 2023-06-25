package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/redis/go-redis/v9"
	"math"
	"strconv"
	"strings"
	"time"
	"verve_challenge/pkg/config"
)

type Item struct {
	Id             string    `json:"id"`
	Price          float64   `json:"price"`
	ExpirationDate time.Time `json:"expiration_date"`
}

func (item *Item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Id             string  `json:"id"`
		Price          float64 `json:"price"`
		ExpirationDate string  `json:"expiration_date"`
	}{
		Id:             item.Id,
		Price:          math.Ceil(item.Price*100) / 100,
		ExpirationDate: item.ExpirationDate.Format(time.DateTime),
	})
}

func ParseItem(s string) (Item, error) {
	item := Item{}
	var parts = strings.Split(s, ",")
	id := parts[0]

	price, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return item, err
	}

	date, err := dateparse.ParseAny(parts[2])
	if err != nil {
		return item, err
	}
	item.ExpirationDate = date
	item.Price = price
	item.Id = id

	return item, nil
}

type Storage struct {
	client *redis.Client
	ctx    context.Context
}

func (s *Storage) AddItem(i Item) {
	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	item := string(b)
	_, err = s.client.Set(s.ctx, i.Id, item, 0).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *Storage) Size() int {
	return int(s.client.DBSize(s.ctx).Val())
}

func (s *Storage) Get(id string) *Item {
	js, err := s.client.Get(s.ctx, id).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var item Item
	err = json.Unmarshal([]byte(js), &item)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &item
}

func (s *Storage) Clean() {
	_, err := s.client.FlushAll(s.ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
}

func CreateStorage(cfg config.Config) (*Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Storage{client: client, ctx: ctx}, nil
}
