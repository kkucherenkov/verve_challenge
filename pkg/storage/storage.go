package storage

import (
	"encoding/json"
	"github.com/araddon/dateparse"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
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
	DB       map[string]Item
	mapMutex sync.RWMutex
}

func (s *Storage) AddItem(i Item) {
	s.mapMutex.Lock()
	s.DB[i.Id] = i
	s.mapMutex.Unlock()
}

func (s *Storage) Size() int {
	return len(s.DB)
}

func (s *Storage) Get(id string) *Item {
	s.mapMutex.RLock()
	item, presence := s.DB[id]
	s.mapMutex.RUnlock()
	if !presence {
		return nil
	}

	return &item
}

func (s *Storage) Clean() {
	s.mapMutex.Lock()
	s.DB = make(map[string]Item)
	s.mapMutex.Unlock()
}

func CreateStorage() *Storage {
	return &Storage{DB: make(map[string]Item)}
}
