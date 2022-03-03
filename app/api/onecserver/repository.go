package onecserver

import "context"

type (
	Repository interface {
		Add(ctx context.Context, u Server) (int64, error)
		Get(ctx context.Context, ID int64) (Server, error)
		GetList(ctx context.Context) ([]Server, error)
		Update(ctx context.Context, u Server) error
		UpdateHash(ctx context.Context, ID int64, hash string) error
		UpdateWatch(ctx context.Context, ID int64, watch bool) error
		Delete(ctx context.Context, ID int64) error
	}
)
