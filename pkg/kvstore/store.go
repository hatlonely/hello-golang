package kvstore

import (
	"github.com/hatlonely/hello-golang/pkg/refx"
)

type Store interface {
	Set(key any, value any) error
	Get(key any) (any, error)
}

func NewStore(options *refx.TypeOptions) (Store, error) {
	kvstore, err := refx.New(options)
	if err != nil {
		return nil, err
	}
	return kvstore.(Store), nil
}
