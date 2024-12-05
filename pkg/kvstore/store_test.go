package kvstore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	_ "github.com/hatlonely/hello-golang/pkg/kvstore/loader"
	_ "github.com/hatlonely/hello-golang/pkg/kvstore/marshaler"
	_ "github.com/hatlonely/hello-golang/pkg/kvstore/store"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

func TestStore(t *testing.T) {
	config := `{
	"Type": "store.Freecache",
	"Options": {
		"Size": 1000000,
		"KeyMarshaler": {
			"Type": "marshaler.JSONMarshaler"
		},
		"ValMarshaler": {
			"Type": "marshaler.JSONMarshaler"
		}
	}
}`

	options := &refx.Options{}
	err := json.Unmarshal([]byte(config), options)
	if err != nil {
		panic(err)
	}

	store, err := kvstore.NewStore(options)
	if err != nil {
		panic(err)
	}

	store.Set(context.Background(), "key", "value")
	val, err := store.Get(context.Background(), "key")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}

func TestStoreGroup(t *testing.T) {
	config := `{
	"Type": "store.StoreGroup",
	"Options": {
		"Stores": [
			{
				"Type": "store.Freecache",
				"Options": {
					"Size": 1000000,
					"KeyMarshaler": {
						"Type": "marshaler.JSONMarshaler"
					},
					"ValMarshaler": {
						"Type": "marshaler.JSONMarshaler"
					}
				}
			},
			{
				"Type": "store.Freecache",
				"Options": {
					"Size": 1000000,
					"KeyMarshaler": {
						"Type": "marshaler.JSONMarshaler"
					},
					"ValMarshaler": {
						"Type": "marshaler.JSONMarshaler"
					}
				}
			}
		]
	}
}`

	options := &refx.Options{}
	err := json.Unmarshal([]byte(config), options)
	if err != nil {
		panic(err)
	}

	store, err := kvstore.NewStore(options)
	if err != nil {
		panic(err)
	}

	store.Set(context.Background(), "key", "value")
	val, err := store.Get(context.Background(), "key")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}

func TestLoadableStore(t *testing.T) {
	config := `{
	"Type": "store.LoadableStore",
	"Options": {
		"Loader": {
			"Type": "loader.FileLoader",
			"Options": {
				"Path": "test.txt"
			}
		},
		"Store": {
			"Type": "store.Freecache",
			"Options": {
				"Size": 1000000,
				"KeyMarshaler": {
					"Type": "marshaler.JSONMarshaler"
				},
				"ValMarshaler": {
					"Type": "marshaler.JSONMarshaler"
				}
			}
		}
	}
}`

	options := &refx.Options{}
	err := json.Unmarshal([]byte(config), options)
	if err != nil {
		panic(err)
	}

	store, err := kvstore.NewStore(options)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)

	// store.Set(context.Background(), "key", "value")
	val, err := store.Get(context.Background(), "key1")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}

type imap map[int]any

type amap map[any]any

func (m amap) Set(key any, value any) {
	m[key] = value
}

func (m amap) Get(key any) any {
	return m[key]
}

func (m imap) Set(key int, value any) {
	m[key] = value
}

func (m imap) Get(key int) any {
	return m[key]
}

type imapStore struct {
	imap imap
}

func (s *imapStore) Set(ctx context.Context, key any, value any) error {
	s.imap[key.(int)] = value
	return nil
}

func (s *imapStore) Get(ctx context.Context, key any) (any, error) {
	return s.imap[key.(int)], nil
}

func BenchmarkMap(b *testing.B) {
	m := imap{}
	for i := 0; i < 1000000; i++ {
		m[i] = i
	}

	b.Run("imap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m.Get(i)
		}
	})

	s := &imapStore{imap: imap{}}
	for i := 0; i < 1000000; i++ {
		s.imap[i] = i
	}

	b.Run("imapStore", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = s.Get(context.Background(), i)
		}
	})

	a := amap{}
	for i := 0; i < 1000000; i++ {
		a[i] = i
	}
	a["hello"] = "world"

	b.Run("amap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = a.Get(i)
		}
	})

	b.Run("none", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// _ = a[i]
		}
	})
}
