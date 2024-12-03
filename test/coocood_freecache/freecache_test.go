package coocoodfreecache_test

import (
	"testing"

	"github.com/coocood/freecache"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoocoodFreecache(t *testing.T) {
	Convey("TestCoocoodFreecache", t, func() {
		c := freecache.NewCache(100 * 1024 * 1024) // 100MB cache

		c.Set([]byte("key"), []byte("val"), 0)
		val, err := c.Get([]byte("key"))
		So(err, ShouldBeNil)
		So(string(val), ShouldEqual, "val")

		c.Set([]byte("int"), []byte("1"), 0)
		i, err := c.Get([]byte("int"))
		So(err, ShouldBeNil)
		So(string(i), ShouldEqual, "1")
	})
}
