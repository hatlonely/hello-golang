package store

import (
	"github.com/coocood/freecache"
	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
	"github.com/pkg/errors"
)

type FreecacheOptions struct {
	Size         int
	KeyMarshaler refx.TypeOptions
	ValMarshaler refx.TypeOptions
}

type Freecache struct {
	cache        *freecache.Cache
	keyMarshaler kvstore.Marshaler
	valMarshaler kvstore.Marshaler
}

func NewFreecache(options FreecacheOptions) (*Freecache, error) {
	keyMarshaler, err := kvstore.NewMarshaler(options.KeyMarshaler)
	if err != nil {
		return nil, errors.Wrap(err, "NewMarshaler failed")
	}

	valMarshaler, err := kvstore.NewMarshaler(options.ValMarshaler)
	if err != nil {
		return nil, errors.Wrap(err, "NewMarshaler failed")
	}

	return &Freecache{
		cache:        freecache.NewCache(options.Size),
		keyMarshaler: keyMarshaler,
		valMarshaler: valMarshaler,
	}, nil
}

func (f *Freecache) Set(key any, value any) error {
	keyBytes, err := f.keyMarshaler.Marshal(key)
	if err != nil {
		return errors.Wrap(err, "Marshal failed")
	}

	valBytes, err := f.valMarshaler.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "Marshal failed")
	}

	f.cache.Set(keyBytes, valBytes, 0)
	return nil
}

func (f *Freecache) Get(key any) (any, error) {
	keyBytes, err := f.keyMarshaler.Marshal(key)
	if err != nil {
		return nil, errors.Wrap(err, "Marshal failed")
	}

	valBytes, err := f.cache.Get(keyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "Get failed")
	}

	var val any
	if err := f.valMarshaler.Unmarshal(valBytes, &val); err != nil {
		return nil, errors.Wrap(err, "Unmarshal failed")
	}

	return val, nil
}
