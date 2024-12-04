package store

import (
	"context"
	"fmt"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

func init() {
	refx.Register("store", "LoadableStore", NewLoadableStore)
}

type LoadableStoreOptions struct {
	Loader refx.Options
	Store  refx.Options
}

type LoadableStore struct {
	loader       kvstore.Loader
	store        kvstore.Store
	storeOptions refx.Options
}

func NewLoadableStore(options *LoadableStoreOptions) (*LoadableStore, error) {
	loader, err := kvstore.NewLoader(&options.Loader)
	if err != nil {
		return nil, err
	}

	store, err := kvstore.NewStore(&options.Store)
	if err != nil {
		return nil, err
	}

	ls := &LoadableStore{
		loader:       loader,
		store:        store,
		storeOptions: options.Store,
	}

	loader.OnChange(func(stream kvstore.KVStream) error {
		store, err := kvstore.NewStore(&options.Store)
		if err != nil {
			return err
		}

		for stream.HasNext() {
			key, val, err := stream.Next()
			fmt.Println(key, val)
			if err != nil {
				return err
			}
			store.Set(context.Background(), key, val)
		}

		ls.store = store

		return nil
	})

	return ls, nil
}

func (s *LoadableStore) Set(ctx context.Context, key any, value any) error {
	return s.store.Set(ctx, key, value)
}

func (s *LoadableStore) Get(ctx context.Context, key any) (any, error) {
	return s.store.Get(ctx, key)
}
