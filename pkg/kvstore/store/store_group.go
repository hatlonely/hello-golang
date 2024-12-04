package store

import (
	"context"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

func init() {
	refx.Register("store", "StoreGroup", NewStoreGroup)
}

type StoreGroupOptions struct {
	Stores []refx.Options
}

type StoreGroup struct {
	stores []kvstore.Store
}

func NewStoreGroup(options *StoreGroupOptions) (*StoreGroup, error) {
	stores := make([]kvstore.Store, 0)
	for _, storeOptions := range options.Stores {
		store, err := kvstore.NewStore(&storeOptions)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	return &StoreGroup{
		stores: stores,
	}, nil
}

func (g *StoreGroup) Set(ctx context.Context, key any, val any) error {
	for _, store := range g.stores {
		err := store.Set(ctx, key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *StoreGroup) Get(ctx context.Context, key any) (any, error) {
	for i, store := range g.stores {
		val, err := store.Get(ctx, key)
		if err == nil {
			for j := 0; j < i; j++ {
				g.stores[j].Set(ctx, key, val)
			}
			return val, nil
		}
	}

	return nil, kvstore.ErrNotFound
}
