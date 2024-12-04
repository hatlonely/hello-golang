package refx

import (
	"encoding/json"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewConstructorInfo(t *testing.T) {
	Convey("TestNewConstructorInfo", t, func() {
		type A interface{}
		type Options struct {
			Key1 string
		}

		Convey("实例", func() {
			constructor, err := NewConstructor("hello world")
			So(err, ShouldBeNil)
			So(constructor, ShouldNotBeNil)
			So(constructor.IsInstance, ShouldBeTrue)
			So(constructor.Instance.Interface(), ShouldEqual, "hello world")
		})

		Convey("没有参数", func() {
			constructor, err := NewConstructor(func() A {
				return "hello world"
			})
			So(err, ShouldBeNil)
			So(constructor, ShouldNotBeNil)
			So(constructor.FuncValue, ShouldNotBeNil)
			So(constructor.IsInstance, ShouldBeFalse)
			So(constructor.HasParam, ShouldBeFalse)
			So(constructor.ReturnError, ShouldBeFalse)
		})

		Convey("一个参数", func() {
			constructor, err := NewConstructor(func(options *Options) A {
				return "hello world"
			})
			So(err, ShouldBeNil)
			So(constructor, ShouldNotBeNil)
			So(constructor.FuncValue, ShouldNotBeNil)
			So(constructor.IsInstance, ShouldBeFalse)
			So(constructor.HasParam, ShouldBeTrue)
			So(constructor.ParamType, ShouldEqual, reflect.TypeOf(&Options{}))
			So(constructor.ReturnError, ShouldBeFalse)
		})

		Convey("一个参数，返回错误", func() {
			constructor, err := NewConstructor(func(options *Options) (A, error) {
				return "hello world", nil
			})
			So(err, ShouldBeNil)
			So(constructor, ShouldNotBeNil)
			So(constructor.FuncValue, ShouldNotBeNil)
			So(constructor.IsInstance, ShouldBeFalse)
			So(constructor.HasParam, ShouldBeTrue)
			So(constructor.ParamType, ShouldEqual, reflect.TypeOf(&Options{}))
			So(constructor.ReturnError, ShouldBeTrue)
		})

		Convey("错误", func() {
			Convey("参数超过两个", func() {
				constructor, err := NewConstructor(func(options *Options, a string, b string) A {
					return "hello world"
				})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "constructor parameters number should not greater than 1")
				So(constructor, ShouldBeNil)
			})
			Convey("第二个参数不是 refx.Option", func() {
				constructor, err := NewConstructor(func(options *Options, a string) A {
					return "hello world"
				})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "constructor parameters number should not greater than 1")
				So(constructor, ShouldBeNil)
			})
			Convey("没有返回值", func() {
				constructor, err := NewConstructor(func(options *Options) {})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "constructor should return value")
				So(constructor, ShouldBeNil)
			})
			Convey("没有返回错误", func() {
				constructor, err := NewConstructor(func(options *Options) (A, string) {
					return "hello world", ""
				})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "constructor second return value should be an error")
				So(constructor, ShouldBeNil)
			})
			Convey("返回值超过两个", func() {
				constructor, err := NewConstructor(func(options *Options) (A, string, string) {
					return "hello world", "", ""
				})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "constructor return number should not greater than 2")
				So(constructor, ShouldBeNil)
			})
		})
	})
}

type B interface {
	Fun() int
}

type C1 struct{}
type C1Options struct{}

func NewC1() (*C1, error) {
	return &C1{}, nil
}

func NewC1WithOptions(options *C1Options) (*C1, error) {
	return &C1{}, nil
}

func (c *C1) Fun() int {
	return 1
}

type C2 struct {
	Options *C2Options
}

type C2Options struct {
	Val int
}

func NewC2WithOptions(options *C2Options) (*C2, error) {
	return &C2{
		Options: options,
	}, nil
}
func (c *C2) Fun() int {
	return c.Options.Val
}

func TestNew(t *testing.T) {
	Convey("TestNew", t, func() {
		Register("", "DefaultC1", NewC1)
		Register("", "C1", NewC1WithOptions)
		Register("", "C2", NewC2WithOptions)
		Register("", "C3", &C2{Options: &C2Options{Val: 3}})
		Register("", "C4", nil)

		{
			v, err := New(&TypeOptions{
				Type: "DefaultC1",
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 1)
		}
		{
			v, err := New(&TypeOptions{
				Type:    "C2",
				Options: json.RawMessage(`{"Val":2}`),
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 2)
		}
		{
			v, err := New(&TypeOptions{
				Type: "C3",
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 3)
		}

		{
			v, err := NewType(reflect.TypeOf((*B)(nil)).Elem(), &TypeOptions{
				Type: "C3",
			})
			So(err, ShouldBeNil)
			So(v.(B).Fun(), ShouldEqual, 3)
		}
		{
			v, err := NewType(reflect.TypeOf((*B)(nil)).Elem(), &TypeOptions{
				Type: "C4",
			})
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
			// 不支持 interface 的 nil 直接转换
			// So(v.(B), ShouldBeNil)
		}

		{
			v, err := NewType(reflect.TypeOf((*C1)(nil)), &TypeOptions{
				Type: "C4",
			})
			So(err, ShouldBeNil)
			So(v, ShouldBeNil)
			So(v.(*C1), ShouldBeNil)
			So(v.(B), ShouldBeNil)
		}
	})
}
