package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	mdl "verve_challenge_web_api/model"
	"verve_challenge_web_api/pkg/model"
	web_api "verve_challenge_web_api/service"
)

func MakeGetItemEndpoint(svc web_api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(mdl.GetItemRequest)
		err, item := svc.GetItem(req.Id)
		if err != nil {
			return model.Item{}, err
		}
		return item, nil
	}
}

func MakeReloadEndpoint(svc web_api.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(mdl.ReloadRequest)
		err := svc.Reload(req.Path)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}
