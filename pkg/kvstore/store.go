package kvstore

import (
	"context"

	"github.com/hatlonely/hello-golang/pkg/refx"
)

type Store interface {
	Set(ctx context.Context, key any, value any) error
	Get(ctx context.Context, key any) (any, error)
	// Del(ctx context.Context, key any) error
}

func NewStore(options *refx.Options) (Store, error) {
	kvstore, err := refx.New(options)
	if err != nil {
		return nil, err
	}
	return kvstore.(Store), nil
}
