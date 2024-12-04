package kvstore

import (
	"github.com/hatlonely/hello-golang/pkg/refx"
)

type KVStore interface {
	Set(key any, value any) error
	Get(key any) (any, error)
}

func NewKVStore(options refx.TypeOptions) (KVStore, error) {
	// 根据配置创建一个KVStore
	kvstore, err := refx.New(&options)
	if err != nil {
		return nil, err
	}
	return kvstore.(KVStore), nil
}

type Marshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
