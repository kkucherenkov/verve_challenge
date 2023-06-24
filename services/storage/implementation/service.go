package implementation

import (
	"context"
	"errors"
	"github.com/araddon/dateparse"
	"strconv"
	"strings"
	"sync"
	"verve_challenge_storage/pkg/file_processor"
	"verve_challenge_storage/pkg/model"
	"verve_challenge_storage/service"
)

type storeSvc struct {
	db       map[string]model.Item
	mapMutex sync.RWMutex
}

func (s *storeSvc) AddItem(line string) error {
	err, item := parseItem(line)
	if err != nil {
		return err
	}
	s.mapMutex.Lock()
	s.db[item.Id] = item
	s.mapMutex.Unlock()
	return nil
}

func (s *storeSvc) Size(_ context.Context) (error, int) {
	return nil, len(s.db)
}

func (s *storeSvc) Get(_ context.Context, id string) (error, model.Item) {
	s.mapMutex.RLock()
	item, presence := s.db[id]
	s.mapMutex.RUnlock()
	if !presence {
		return errors.New("not found"), model.Item{}
	}

	return nil, item
}

func (s *storeSvc) clean(context.Context) error {
	s.mapMutex.Lock()
	s.db = make(map[string]model.Item)
	s.mapMutex.Unlock()
	return nil
}

func (s *storeSvc) Reload(ctx context.Context, path string) error {
	err := s.clean(ctx)
	if err != nil {
		return err
	}
	err = file_processor.ProcessFile(path, s)
	if err != nil {
		return err
	}
	return nil
}

func parseItem(s string) (error, model.Item) {
	item := model.Item{}
	var parts = strings.Split(s, ",")
	id := parts[0]

	price, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return err, item
	}

	date, err := dateparse.ParseAny(parts[2])
	if err != nil {
		return err, item
	}
	item.ExpirationDate = date
	item.Price = float32(price)
	item.Id = id

	return nil, item
}

func CreateStorage() service.Service {
	return &storeSvc{
		db: make(map[string]model.Item),
	}
}
