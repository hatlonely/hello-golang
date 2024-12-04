package kvstore

import (
	"github.com/hatlonely/hello-golang/pkg/refx"
)

type KVStore interface {
	Set(key any, value any) error
	Get(key any) (any, error)
}

func NewKVStore(options refx.TypeOptions) (KVStore, error) {
	kvstore, err := refx.New(&options)
	if err != nil {
		return nil, err
	}
	return kvstore.(KVStore), nil
}
