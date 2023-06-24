package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"math"
	"time"
	"verve_challenge_web_api/client/service"
	"verve_challenge_web_api/pkg/model"
)

type GetSizeRequest struct {
}

type GetSizeResponse struct {
	Size int64  `json:"size"`
	Err  string `json:"error,omitempty"`
}

type GetItemRequest struct {
	Id string `json:"id"`
}

type GetItemResponse struct {
	Id             string  `json:"id"`
	Price          float32 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
	Err            string  `json:"err,omitempty"`
}

type ReloadDbRequest struct {
	Path string `json:"path"`
}

type ReloadDbResponse struct {
	Err string `json:"err,omitempty"`
}

type Endpoints struct {
	GetSizeEndpoint  endpoint.Endpoint
	GetItemEndpoint  endpoint.Endpoint
	ReloadDbEndpoint endpoint.Endpoint
}

func MakeGetSizeEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		err, size := svc.Size(ctx)
		if err != nil {
			return nil, err
		}
		return GetSizeResponse{Size: int64(size)}, nil
	}
}

func MakeGetItemEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetItemRequest)
		err, item := svc.Get(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		return GetItemResponse{Id: item.Id, Price: float32(math.Ceil(float64(item.Price*100)) / 100), ExpirationDate: item.ExpirationDate.Format(time.DateTime)}, nil
	}
}

func MakeReloadDbEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ReloadDbRequest)
		err = svc.Reload(ctx, req.Path)
		if err != nil {
			return nil, err
		}
		return ReloadDbResponse{}, nil
	}
}

func (e Endpoints) Size(ctx context.Context) (error, int) {
	req := GetSizeRequest{}
	resp, err := e.GetSizeEndpoint(ctx, req)
	if err != nil {
		return err, -1
	}
	getSizeResponse := resp.(GetSizeResponse)
	if getSizeResponse.Err != "" {
		return err, -1
	}
	return nil, int(getSizeResponse.Size)
}

func (e Endpoints) Get(ctx context.Context, id string) (error, model.Item) {
	req := GetItemRequest{Id: id}
	resp, err := e.GetItemEndpoint(ctx, req)
	if err != nil {
		return err, model.Item{}
	}
	getItemResponse := resp.(GetItemResponse)
	if getItemResponse.Err != "" {
		return err, model.Item{}
	}
	parsedDate, err := time.Parse(time.DateTime, getItemResponse.ExpirationDate)
	if err != nil {
		return err, model.Item{}
	}
	item := model.Item{
		Id:             getItemResponse.Id,
		Price:          getItemResponse.Price,
		ExpirationDate: parsedDate,
	}
	return nil, item
}

func (e Endpoints) Reload(ctx context.Context, path string) error {
	req := ReloadDbRequest{Path: path}
	_, err := e.ReloadDbEndpoint(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
