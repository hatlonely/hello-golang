package store

import (
	"context"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

func init() {
	refx.Register("store", "LoadableStore", NewLoadableStore)
}

type LoadableStoreOptions struct {
	Store         refx.Options
	ReadonlyStore refx.Options
}

type LoadableStore struct {
	store         kvstore.Store
	readonlyStore kvstore.ReadonlyStore
}

func NewLoadableStore(options LoadableStoreOptions) (*LoadableStore, error) {
	store, err := kvstore.NewStore(&options.Store)
	if err != nil {
		return nil, err
	}

	readonlyStore, err := kvstore.NewReadonlyStore(&options.ReadonlyStore)
	if err != nil {
		return nil, err
	}

	return &LoadableStore{
		store:         store,
		readonlyStore: readonlyStore,
	}, nil
}

func (s *LoadableStore) Set(ctx context.Context, key any, value any) error {
	return s.store.Set(ctx, key, value)
}

func (s *LoadableStore) Get(ctx context.Context, key any) (any, error) {
	val, err := s.store.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	if err != kvstore.ErrNotFound {
		return nil, err
	}

	val, err = s.readonlyStore.Get(ctx, key)
	if err == nil {
		s.store.Set(ctx, key, val)
		return val, nil
	}

	if err == kvstore.ErrNotFound {
		s.store.Set(ctx, key, nil)
		return nil, kvstore.ErrNotFound
	}

	return nil, err
}
