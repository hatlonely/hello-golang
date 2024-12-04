package store

import (
	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

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
