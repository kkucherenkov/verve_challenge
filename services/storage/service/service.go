package service

import (
	"context"
	"verve_challenge_storage/pkg/model"
)

type Service interface {
	Size(ctx context.Context) (error, int)
	Get(ctx context.Context, id string) (error, model.Item)
	Reload(ctx context.Context, path string) error
}
