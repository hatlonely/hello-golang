package kvstore_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	_ "github.com/hatlonely/hello-golang/pkg/kvstore/marshaler"
	_ "github.com/hatlonely/hello-golang/pkg/kvstore/store"
	"github.com/hatlonely/hello-golang/pkg/refx"
)

func TestStore(t *testing.T) {
	config := `{
	"Namespace": "store",
	"Type": "Freecache",
	"Options": {
		"Size": 1000000,
		"KeyMarshaler": {
			"Namespace": "marshaler",
			"Type": "JSONMarshaler"
		},
		"ValMarshaler": {
			"Namespace": "marshaler",
			"Type": "JSONMarshaler"
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

	store.Set("key", "value")
	val, err := store.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
