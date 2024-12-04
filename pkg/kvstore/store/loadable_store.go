package store

import (
	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

type LoadableStoreOptions struct {
	Loader refx.Options
	Store  refx.Options
}

type LoadableStore struct {
	loader kvstore.Loader
	store  kvstore.Store
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

	return &LoadableStore{
		loader: loader,
		store:  store,
	}, nil
}
