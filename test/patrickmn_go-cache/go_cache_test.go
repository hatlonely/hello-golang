package patrickmngocache_test

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGoCache(t *testing.T) {
	Convey("TestGoCache", t, func() {
		c := cache.New(5*time.Minute, 10*time.Minute)

		c.Set("key", "val", cache.DefaultExpiration)
		val, ok := c.Get("key")
		So(ok, ShouldBeTrue)
		So(val, ShouldEqual, "val")

		c.Set("int", 1, cache.DefaultExpiration)
		i, ok := c.Get("int")
		So(ok, ShouldBeTrue)
		So(i, ShouldEqual, 1)
	})
}
