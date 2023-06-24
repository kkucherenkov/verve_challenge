package model

import (
	"context"
	"verve_challenge_storage/endpoints"
	"verve_challenge_storage/pb"
)

func EncodeGRPCReloadResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.ReloadDbResponse)
	return &pb.ReloadDbResponse{Err: req.Err}, nil
}

func DecodeGRPCReloadResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ReloadDbResponse)
	return endpoints.ReloadDbResponse{
		Err: req.Err,
	}, nil
}

func EncodeGRPCReloadRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.ReloadDbRequest)
	return &pb.ReloadDbRequest{Path: req.Path}, nil
}

func DecodeGRPCReloadRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ReloadDbRequest)
	return endpoints.ReloadDbRequest{
		Path: req.Path,
	}, nil
}

func EncodeGRPCGetSizeResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.GetSizeResponse)
	return &pb.GetSizeResponse{
		Err:  req.Err,
		Size: req.Size,
	}, nil
}

func DecodeGRPCGetSizeResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetSizeResponse)
	return endpoints.GetSizeResponse{
		Err:  req.Err,
		Size: req.Size,
	}, nil
}

func EncodeGRPCGetSizeRequest(_ context.Context, r interface{}) (interface{}, error) {
	return &pb.GetSizeRequest{}, nil
}

func DecodeGRPCGetSizeRequest(_ context.Context, r interface{}) (interface{}, error) {
	return endpoints.GetSizeRequest{}, nil
}

func EncodeGRPCGetItemResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.GetItemResponse)
	return &pb.GetItemResponse{
		Err:            req.Err,
		Id:             req.Id,
		Price:          req.Price,
		ExpirationDate: req.ExpirationDate,
	}, nil
}

func DecodeGRPCGetItemResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetItemResponse)
	return endpoints.GetItemResponse{
		Err:            req.Err,
		Id:             req.Id,
		Price:          req.Price,
		ExpirationDate: req.ExpirationDate,
	}, nil
}

func EncodeGRPCGetItemRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetItemRequest)
	return &endpoints.GetItemRequest{Id: req.Id}, nil
}

func DecodeGRPCGetItemRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.GetItemRequest)
	return endpoints.GetItemRequest{
		Id: req.Id,
	}, nil
}
