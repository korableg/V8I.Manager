package onecserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
	"io"
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
		io.Closer
	}

	service struct {
		serviceRepo Repository
		dbCollector onecdb.DBCollector

		mu         sync.Mutex
		watching   map[int64]chan<- struct{}
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
		watching:    make(map[int64]chan<- struct{}, 0),
	}

	if err := s.startWatchingAllServers(); err != nil {
		return nil, errors.New("db collector is not defined")
	}

	return s, nil
}

func (s *service) Add(ctx context.Context, req AddServerRequest) (int64, error) {
	server := Server{
		Name:    req.Name,
		LSTPath: req.LSTPath,
	}

	id, err := s.serviceRepo.Add(ctx, server)
	if err != nil {
		return 0, fmt.Errorf("add to store: %w", err)
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, ID int64) (Server, error) {
	server, err := s.serviceRepo.Get(ctx, ID)
	if err != nil {
		return Server{}, fmt.Errorf("get from store: %w", err)
	}

	s.mu.Lock()
	_, server.Watch = s.watching[ID]
	s.mu.Unlock()

	return server, nil
}

func (s *service) GetList(ctx context.Context) ([]Server, error) {
	servers, err := s.serviceRepo.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get from store: %w", err)
	}

	s.mu.Lock()
	for _, server := range servers {
		_, server.Watch = s.watching[server.ID]
	}
	s.mu.Unlock()

	return servers, nil
}

func (s *service) Update(ctx context.Context, req UpdateServerRequest) error {
	if err := s.serviceRepo.Update(ctx, Server{
		ID:      req.ID,
		Name:    req.Name,
		LSTPath: req.LSTPath,
	}); err != nil {
		return fmt.Errorf("update server in store: %w", err)
	}

	return nil
}

func (s *service) SwitchWatching(ctx context.Context, ID int64) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stopCh, ok := s.watching[ID]
	if ok {
		close(stopCh)
		delete(s.watching, ID)
	} else {
		stopCh, err = s.startWatching(ID)
		if err != nil {
			return fmt.Errorf("switch watching: %w", err)
		}
		s.watching[ID] = stopCh
	}

	if err = s.serviceRepo.UpdateWatch(ctx, ID, !ok); err != nil {
		return fmt.Errorf("update watching in storage: %w", err)
	}

	return nil
}

func (s *service) Delete(ctx context.Context, ID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, ok := s.watching[ID]; ok {
		close(ch)
		delete(s.watching, ID)
	}

	if err := s.serviceRepo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete server from store: %w", err)
	}

	return nil
}

func (s *service) Close() error {
	s.mu.Lock()
	for _, ch := range s.watching {
		close(ch)
	}
	s.watching = make(map[int64]chan<- struct{}, 0)
	s.mu.Unlock()

	s.watchingWg.Wait()
	return nil
}

func (s *service) startWatching(id int64) (chan<- struct{}, error) {
	ch := make(chan struct{}, 1)

	return ch, nil
}

func (s *service) startWatchingAllServers() error {
	return nil
}
