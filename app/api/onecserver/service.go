package onecserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/korableg/V8I.Manager/app/api/onecdb"
	"github.com/korableg/V8I.Manager/app/api/onecserver/dbbuilder"
	"github.com/korableg/V8I.Manager/app/api/onecserver/watcher"
	"github.com/sirupsen/logrus"
	"sync"
)

type (
	service struct {
		serviceRepo   Repository
		dbCollector   onecdb.DBCollector
		watcherFabric watcher.Fabric
		builder       dbbuilder.DBBuilder

		mu         sync.Mutex
		watching   map[int64]chan<- struct{}
		watchingWg sync.WaitGroup
	}
)

func NewService(serviceRepo Repository, collector onecdb.DBCollector, watcherFabric watcher.Fabric, builder dbbuilder.DBBuilder) (*service, error) {
	if serviceRepo == nil {
		return nil, errors.New("service repository is not defined")
	}

	if collector == nil {
		return nil, errors.New("db collector is not defined")
	}

	if watcherFabric == nil {
		return nil, errors.New("watcher fabric is not defined")
	}

	if builder == nil {
		return nil, errors.New("db builder is not defined")
	}

	s := &service{
		serviceRepo:   serviceRepo,
		dbCollector:   collector,
		watcherFabric: watcherFabric,
		builder:       builder,
		watching:      make(map[int64]chan<- struct{}, 0),
	}

	if err := s.startWatchingAllServers(); err != nil {
		return nil, errors.New("start watching all servers")
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

func (s *service) SwitchWatching(ctx context.Context, ID int64) (watching bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stopCh, ok := s.watching[ID]
	if ok {
		close(stopCh)
		delete(s.watching, ID)
	} else {
		stopCh, err = s.startWatching(ctx, ID)
		if err != nil {
			return false, fmt.Errorf("switch watching: %w", err)
		}
		s.watching[ID] = stopCh
	}

	if err = s.serviceRepo.UpdateWatch(ctx, ID, !ok); err != nil {
		logrus.Errorf("update watch in repository: %s", err.Error())
	}

	return !ok, nil
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

func (s *service) startWatching(ctx context.Context, id int64) (chan<- struct{}, error) {
	server, err := s.serviceRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get server from repository: %w", err)
	}

	return s.startWatchingServer(server)
}

func (s *service) startWatchingAllServers() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	servers, err := s.serviceRepo.GetList(context.Background())
	if err != nil {
		return fmt.Errorf("get servers from repository: %w", err)
	}

	for _, server := range servers {
		if !server.Watch {
			continue
		}

		stopCh, err := s.startWatchingServer(server)
		if err != nil {
			return fmt.Errorf("start watching server: %w", err)
		}

		s.watching[server.ID] = stopCh
	}

	return nil
}

func (s *service) startWatchingServer(server Server) (chan<- struct{}, error) {
	stopCh := make(chan struct{}, 1)

	w, err := s.watcherFabric(server.LSTPath, stopCh)
	if err != nil {
		return nil, fmt.Errorf("init watcher: %w", err)
	}

	s.watchingWg.Add(1)
	go func() {
		defer s.watchingWg.Done()

		for range w.Start() {
			dbs, err := s.builder.Build(server.LSTPath)
			if err != nil {
				logrus.Errorf("build list of db: %s", err.Error())
				continue
			}
			if err = s.dbCollector.Collect(dbs...); err != nil {
				logrus.Errorf("collect db: %s", err.Error())
				continue
			}
		}
	}()

	return stopCh, err
}
