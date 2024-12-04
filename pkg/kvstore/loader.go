package kvstore

import "github.com/hatlonely/hello-golang/pkg/refx"

type KVStream interface {
	HasNext() bool
	Next() (any, any, error)
}

type Listener func(KVStream) error

type Loader interface {
	OnChange(Listener) error

	Close() error
}

func NewLoader(options *refx.Options) (Loader, error) {
	loader, err := refx.New(options)
	if err != nil {
		return nil, err
	}
	return loader.(Loader), nil
}
