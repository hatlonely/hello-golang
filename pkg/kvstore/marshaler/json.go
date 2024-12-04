package marshaler

import (
	"github.com/hatlonely/hello-golang/pkg/refx"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	refx.Register("marshaler", "JSONMarshaler", &JSONMarshaler{})
}

type JSONMarshaler struct{}

func (m *JSONMarshaler) Marshal(v any) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (m *JSONMarshaler) Unmarshal(data []byte, v any) error {
	return jsoniter.Unmarshal(data, v)
}
