package service

import (
	"context"
	"google.golang.org/grpc"
	"verve_challenge_web_api/client"
	svcClient "verve_challenge_web_api/client/service"
	"verve_challenge_web_api/pkg/model"
)

type Service interface {
	GetItem(id string) (error, model.Item)
	Reload(path string) error
}

type httService struct {
	client svcClient.Service
	ctx    context.Context
}

func (c *httService) GetItem(id string) (error, model.Item) {
	err, item := c.client.Get(c.ctx, id)
	if err != nil {
		return err, model.Item{}
	}
	return nil, item
}

func (c *httService) Reload(path string) error {
	err := c.client.Reload(c.ctx, path)
	if err != nil {
		return err
	}
	return nil
}

func NewService(ctx context.Context, conn *grpc.ClientConn) Service {
	return &httService{
		ctx:    ctx,
		client: client.New(conn),
	}
}
