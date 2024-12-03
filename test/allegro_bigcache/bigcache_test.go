package allegrobigcache

import (
	"context"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAllegroBigcache(t *testing.T) {
	Convey("TestAllegroBigcache", t, func() {
		c, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
		So(err, ShouldBeNil)

		c.Set("key", []byte("val"))
		val, err := c.Get("key")
		So(err, ShouldBeNil)
		So(string(val), ShouldEqual, "val")

		c.Set("int", []byte("1"))
		i, err := c.Get("int")
		So(err, ShouldBeNil)
		So(string(i), ShouldEqual, "1")
	})
}
