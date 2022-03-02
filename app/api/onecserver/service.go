package onecserver

import (
	"context"
	"errors"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
	"sync"
)

type (
	Service interface {
		Add(ctx context.Context, u AddServerRequest) (int64, error)
		Get(ctx context.Context, ID int64) (Server, error)
		GetList(ctx context.Context) ([]Server, error)
		Update(ctx context.Context, u UpdateServerRequest) error
		SwitchWatching(ctx context.Context, ID int64) error
		Delete(ctx context.Context, ID int64) error
	}

	service struct {
		serviceRepo Repository
		dbCollector onecdb.DBCollector

		mu         sync.Mutex
		watching   map[int64]chan bool
		watchingWg sync.WaitGroup
	}
)

func NewService(serviceRepo Repository, collector onecdb.DBCollector) (*service, error) {
	if serviceRepo == nil {
		return nil, errors.New("service repository is not defined")
	}

	if collector == nil {
		return nil, errors.New("db collector is not defined")
	}

	s := &service{
		serviceRepo: serviceRepo,
		dbCollector: collector,
	}

	return s, nil
}

func (s *service) Add(ctx context.Context, u AddServerRequest) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Get(ctx context.Context, ID int64) (Server, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetList(ctx context.Context) ([]Server, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Update(ctx context.Context, u UpdateServerRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) SwitchWatching(ctx context.Context, ID int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) Delete(ctx context.Context, ID int64) error {
	//TODO implement me
	panic("implement me")
}
