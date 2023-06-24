package transport

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
	"verve_challenge_storage/endpoints"
	"verve_challenge_storage/model"
	"verve_challenge_storage/pb"
)

type grpcServer struct {
	getSize grpctransport.Handler
	getItem grpctransport.Handler
	reload  grpctransport.Handler
}

// create new grpc server
func NewGRPCServer(_ context.Context, endpoint endpoints.Endpoints) pb.StorageServer {
	return &grpcServer{
		reload: grpctransport.NewServer(
			endpoint.ReloadDbEndpoint,
			model.DecodeGRPCReloadRequest,
			model.EncodeGRPCReloadResponse),
		getItem: grpctransport.NewServer(
			endpoint.GetItemEndpoint,
			model.DecodeGRPCGetItemRequest,
			model.EncodeGRPCGetItemResponse),
		getSize: grpctransport.NewServer(
			endpoint.GetSizeEndpoint,
			model.DecodeGRPCGetSizeRequest,
			model.EncodeGRPCGetSizeRequest),
	}
}

func (s *grpcServer) Size(ctx context.Context, r *pb.GetSizeRequest) (*pb.GetSizeResponse, error) {
	_, resp, err := s.getSize.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetSizeResponse), nil
}

func (s *grpcServer) GetItem(ctx context.Context, r *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	_, resp, err := s.getItem.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetItemResponse), nil
}

func (s *grpcServer) ReloadDb(ctx context.Context, r *pb.ReloadDbRequest) (*pb.ReloadDbResponse, error) {
	_, resp, err := s.reload.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ReloadDbResponse), nil
}
