package kvstore

type KVStream interface {
	HasNext() bool
	Next() (any, any, error)
}

type Listener func(KVStream) error

type Loader interface {
	OnChange(Listener) error

	Close() error
}
