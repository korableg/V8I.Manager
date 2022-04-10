package client

import "context"

type (
	service struct {
	}
)

func NewService() (*service, error) {
	return &service{}, nil
}

func (s *service) NewClient(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}
