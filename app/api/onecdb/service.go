package onecdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type (
	DBCollector interface {
		Collect(db ...DB) error
	}

	Service interface {
		Add(ctx context.Context, reqDB AddDBRequest) (int64, error)
		Get(ctx context.Context, ID int64) (DB, error)
		GetList(ctx context.Context) ([]DB, error)
		Update(ctx context.Context, reqDB UpdateDBRequest) error
		Delete(ctx context.Context, ID int64) error
	}

	service struct {
		dbRepo Repository
	}
)

func NewService(dbRepo Repository) (*service, error) {
	if dbRepo == nil {
		return nil, errors.New("db repository is not defined")
	}

	s := &service{
		dbRepo: dbRepo,
	}

	return s, nil
}

func (s *service) Add(ctx context.Context, reqDB AddDBRequest) (int64, error) {
	reqUUID, err := uuid.Parse(reqDB.UUID)
	if err != nil {
		return 0, fmt.Errorf("parse UUID: %w", err)
	}

	db := DB{
		UUID:                  reqUUID,
		Name:                  reqDB.Name,
		Connect:               reqDB.Connect,
		OrderInList:           reqDB.OrderInList,
		OrderInTree:           reqDB.OrderInTree,
		Folder:                reqDB.Folder,
		ClientConnectionSpeed: reqDB.ClientConnectionSpeed,
		App:                   reqDB.App,
		WA:                    reqDB.WA,
		Version:               reqDB.Version,
		AdditionalParameters:  reqDB.AdditionalParameters,
	}

	id, err := s.dbRepo.Add(ctx, db)
	if err != nil {
		return 0, fmt.Errorf("add to store: %w", err)
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, ID int64) (DB, error) {
	db, err := s.dbRepo.Get(ctx, ID)
	if err != nil {
		return DB{}, fmt.Errorf("get from store: %w", err)
	}

	return db, nil
}

func (s *service) GetList(ctx context.Context) ([]DB, error) {
	dbs, err := s.dbRepo.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get from store: %w", err)
	}

	return dbs, nil
}

func (s *service) Update(ctx context.Context, reqDB UpdateDBRequest) error {
	reqUUID, err := uuid.Parse(reqDB.UUID)
	if err != nil {
		return fmt.Errorf("parse UUID: %w", err)
	}

	db := DB{
		ID:                    reqDB.ID,
		UUID:                  reqUUID,
		Name:                  reqDB.Name,
		Connect:               reqDB.Connect,
		OrderInList:           reqDB.OrderInList,
		OrderInTree:           reqDB.OrderInTree,
		Folder:                reqDB.Folder,
		ClientConnectionSpeed: reqDB.ClientConnectionSpeed,
		App:                   reqDB.App,
		WA:                    reqDB.WA,
		Version:               reqDB.Version,
		AdditionalParameters:  reqDB.AdditionalParameters,
	}

	if err := s.dbRepo.Update(ctx, db); err != nil {
		return fmt.Errorf("update db in store: %w", err)
	}

	return nil
}

func (s *service) Delete(ctx context.Context, ID int64) error {
	if err := s.dbRepo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete db from store: %w", err)
	}

	return nil
}

func (s *service) Collect(dbs ...DB) error {
	ctx := context.Background()

	for _, db := range dbs {
		if _, err := s.dbRepo.Add(ctx, db); err != nil {
			return fmt.Errorf("add to store: %w", err)
		}
	}

	return nil
}
