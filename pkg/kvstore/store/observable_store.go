package store

import (
	"context"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

type ObservableStoreOptions struct {
	Store refx.Options
	// trace metric log retry

	EnableMetric bool
	EnableTrace  bool
	EnableLog    bool
}

type ObservableStore struct {
	store kvstore.Store
}

func NewObservableStore(options *ObservableStoreOptions) (*ObservableStore, error) {
	store, err := kvstore.NewStore(&options.Store)
	if err != nil {
		return nil, err
	}

	os := &ObservableStore{
		store: store,
	}

	return os, nil
}

func (os *ObservableStore) Set(ctx context.Context, key, val interface{}) error {
	// trace metric log retry

	return os.store.Set(ctx, key, val)
}

func (os *ObservableStore) Get(ctx context.Context, key interface{}) (interface{}, error) {
	return os.store.Get(ctx, key)
}
