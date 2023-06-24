package client

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"verve_challenge_web_api/client/endpoints"
	"verve_challenge_web_api/client/model"
	"verve_challenge_web_api/client/pb"
	"verve_challenge_web_api/client/service"
)

func New(conn *grpc.ClientConn) service.Service {
	var getSizeEndpoint = grpctransport.NewClient(
		conn, "pb.Storage", "GetSize",
		model.EncodeGRPCGetSizeRequest,
		model.DecodeGRPCGetSizeResponse,
		pb.GetSizeResponse{},
	).Endpoint()
	var getItemEndpoint = grpctransport.NewClient(
		conn, "pb.Storage", "GetItem",
		model.EncodeGRPCGetItemRequest,
		model.DecodeGRPCGetItemResponse,
		pb.GetItemResponse{},
	).Endpoint()
	var reloadDbEndpoint = grpctransport.NewClient(
		conn, "pb.Storage", "ReloadDb",
		model.EncodeGRPCReloadRequest,
		model.DecodeGRPCReloadResponse,
		pb.ReloadDbResponse{},
	).Endpoint()

	return endpoints.Endpoints{
		GetItemEndpoint:  getItemEndpoint,
		GetSizeEndpoint:  getSizeEndpoint,
		ReloadDbEndpoint: reloadDbEndpoint,
	}
}
