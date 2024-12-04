package kvstore

import "github.com/hatlonely/hello-golang/pkg/refx"

type Marshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

func NewMarshaler(options refx.Options) (Marshaler, error) {
	marshaler, err := refx.New(&options)
	if err != nil {
		return nil, err
	}
	return marshaler.(Marshaler), nil
}
