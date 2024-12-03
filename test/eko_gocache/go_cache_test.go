package ekogocache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	ristretto_store "github.com/eko/gocache/store/ristretto/v4"
	gocache "github.com/patrickmn/go-cache"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRistrettoCache(t *testing.T) {
	Convey("Given a new Ristretto cache store", t, func() {
		ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1000,
			MaxCost:     100,
			BufferItems: 64,
		})
		So(err, ShouldBeNil)

		myStore := ristretto_store.NewRistretto(ristrettoCache)
		myCache := cache.New[string](myStore)

		err = myCache.Set(context.Background(), "key", "value")
		So(err, ShouldBeNil)

		value, err := myCache.Get(context.Background(), "key")
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "value")
	})
}

func TestRistrettoCache1(t *testing.T) {
	ctx := context.Background()
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1000,
		MaxCost:     100,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)

	cacheManager := cache.New[string](ristrettoStore)
	err = cacheManager.Set(ctx, "my-key", "my-value", store.WithCost(2))
	if err != nil {
		panic(err)
	}

	value, err := cacheManager.Get(ctx, "my-key")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", value)
}

func TestGoCache(t *testing.T) {
	ctx := context.Background()

	gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
	gocacheStore := gocache_store.NewGoCache(gocacheClient)

	cacheManager := cache.New[[]byte](gocacheStore)
	err := cacheManager.Set(ctx, "my-key", []byte("my-value"))
	if err != nil {
		panic(err)
	}

	value, err := cacheManager.Get(ctx, "my-key")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", value)
}
