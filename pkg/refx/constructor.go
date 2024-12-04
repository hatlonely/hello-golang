// 提供一种通用的对象创建方式

package refx

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}

var constructorMap = map[string]*Constructor{}

type Constructor struct {
	IsNil       bool          // 是否为空，一些场景 nil 也是一个合法的值，比如 RateLimiter 为 nil 时，就没有限流逻辑
	HasParam    bool          // 是否包含参数，有些对象可以不需要参数构造
	ReturnError bool          // 返回值中是否包含错误
	IsInstance  bool          // 是否是一个普通对象，而不是一个函数
	Instance    reflect.Value // 普通对象
	ParamType   reflect.Type  // 参数类型
	FuncValue   reflect.Value // 构造函数
}

// 调用构造函数，供 New 来使用
func (c *Constructor) Call(v json.RawMessage) ([]reflect.Value, error) {
	if c.IsNil {
		var v interface{}
		return []reflect.Value{reflect.ValueOf(&v).Elem()}, nil
	}

	// 如果本身注册的是一个实例，直接返回这个实例
	if c.IsInstance {
		return []reflect.Value{c.Instance}, nil
	}

	var params []reflect.Value
	if c.HasParam {
		if reflect.TypeOf(v) == c.ParamType {
			// 如果参数的类型和构造函数的参数类型一致，使用该参数调用函数
			params = append(params, reflect.ValueOf(v))
		} else {
			// 否则使用 InterfaceToStruct 将参数转成构造函数需要的参数类型，从配置中获取的结构会走到这个分支
			param := reflect.New(c.ParamType)
			if err := json.Unmarshal(v, param.Interface()); err != nil {
				return nil, errors.WithMessage(err, "json.Unmarshal failed")
			}
			params = append(params, param.Elem())
		}
	}

	return c.FuncValue.Call(params), nil
}

// 分析构造函数的信息，供 New 来使用
func NewConstructor(constructor interface{}) (*Constructor, error) {
	if constructor == nil {
		return &Constructor{
			IsNil: true,
		}, nil
	}

	rt := reflect.TypeOf(constructor)

	// 注册一个实例，而非一个构造函数，如果构造的对象本身就是一个函数，将不再支持，必须提供一个构造函数
	if rt.Kind() != reflect.Func {
		return &Constructor{
			IsInstance: true,
			Instance:   reflect.ValueOf(constructor),
		}, nil
	}

	var info Constructor
	info.FuncValue = reflect.ValueOf(constructor)
	// 默认没有参数 NewXXX()

	// 输入参数不能超过1个
	if rt.NumIn() > 1 {
		return nil, errors.New("constructor parameters number should not greater than 1")
	}
	// 一个参数 NewXXX(options *Options)
	if rt.NumIn() == 1 {
		info.HasParam = true
	}

	// 返回值不能超过两个
	if rt.NumOut() > 2 {
		return nil, errors.New("constructor return number should not greater than 2")
	}
	// 返回值如果有两个，第二个返回值必须是 error
	if rt.NumOut() == 2 && rt.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, errors.New("constructor second return value should be an error")
	}
	if rt.NumOut() == 0 {
		return nil, errors.New("constructor should return value")
	}

	info.ReturnError = rt.NumOut() == 2

	if info.HasParam {
		info.ParamType = rt.In(0)
	}

	return &info, nil
}

// 通用的配置结构，Options 可以是任意类型，一般来自于配置文件，也可以由代码构造
type TypeOptions struct {
	Namespace string          // 命名空间，一般以 <package>.<interface> 命名
	Type      string          // 类型，一般使用类名或者类名的缩写
	Options   json.RawMessage // 传给构造函数的参数，也可以是一个可以自然转换成参数的 interface{}（一般来自于配置文件）
}

// 注册一个构造方法
func Register(namespace string, typ string, constructor interface{}) {
	key := typ
	if namespace != "" {
		key = fmt.Sprintf("%s.%s", namespace, typ)
	}
	if _, ok := constructorMap[key]; ok {
		panic(fmt.Sprintf("namespace: [%s], type: [%s] is already registered", namespace, typ))
	}

	info, err := NewConstructor(constructor)
	Must(err)

	constructorMap[key] = info
}

// 返回对应构造函数构造的对象，用户需要在自己代码中断言返回的对象是否是自己期望的类型
// v, err := refx.New(options)
//
//	if err != nil {
//	  return nil, errors.WithMessage(err, "refx.New failed")
//	}
//
// i, ok := v.(MyType)
//
//	if !ok {
//	  return nil, errors.New("type assertion failed")
//	}
func New(options *TypeOptions) (interface{}, error) {
	key := options.Type
	if options.Namespace != "" {
		key = fmt.Sprintf("%s.%s", options.Namespace, options.Type)
	}

	if key == "" {
		return nil, nil
	}

	constructor, ok := constructorMap[key]
	if !ok {
		return nil, errors.Errorf("unregistered type: [%v]", key)
	}

	result, err := constructor.Call(options.Options)
	if err != nil {
		return nil, errors.WithMessage(err, "constructor.Call failed")
	}

	if constructor.ReturnError {
		if !result[1].IsNil() {
			return nil, errors.Wrapf(result[1].Interface().(error), "New failed. namespace: [%s], type: [%s]", options.Namespace, options.Type)
		}
		return result[0].Interface(), nil
	}

	return result[0].Interface(), nil
}

// 在 New 的基础上，校验返回值是否是用户期望的类型
// v, err := refx.NewType(reflect.TypeOf((*MyType)(nil)).Elem(), options, opts...)
//
//	if err != nil {
//	  return nil, errors.WithMessage(err, "refx.New failed")
//	}
func NewType(implement reflect.Type, options *TypeOptions) (interface{}, error) {
	v, err := New(options)
	if err != nil {
		return nil, errors.WithMessage(err, "refx.New failed")
	}

	if v == nil {
		return reflect.New(implement).Elem().Interface(), nil
	}

	rv := reflect.TypeOf(v)
	if rv == implement || rv.Implements(implement) {
		return v, nil
	}

	return nil, errors.Errorf("%v is not a %v", rv, implement)
}
