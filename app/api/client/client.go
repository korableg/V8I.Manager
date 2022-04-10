package client

import "context"

type (
	Service interface {
		NewClient(ctx context.Context) (string, error)
	}
)
