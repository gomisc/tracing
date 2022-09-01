package tracing

import (
	"go.opentelemetry.io/otel/attribute"
)

type attributes []attribute.KeyValue

func (a *attributes) Bool(key string, value bool) {
	*a = append(*a, attribute.Bool(key, value))
}

func (a *attributes) Int(key string, value int) {
	*a = append(*a, attribute.Int(key, value))
}

func (a *attributes) Int8(key string, value int8) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Int16(key string, value int16) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Int32(key string, value int32) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Int64(key string, value int64) {
	*a = append(*a, attribute.Int64(key, value))
}

func (a *attributes) Uint(key string, value uint) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Uint8(key string, value uint8) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Uint16(key string, value uint16) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Uint32(key string, value uint32) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Uint64(key string, value uint64) {
	*a = append(*a, attribute.Int(key, int(value)))
}

func (a *attributes) Float32(key string, value float32) {
	*a = append(*a, attribute.Float64(key, float64(value)))
}

func (a *attributes) Float64(key string, value float64) {
	*a = append(*a, attribute.Float64(key, value))
}

func (a *attributes) Str(key string, value string) {
	*a = append(*a, attribute.String(key, value))
}

func (a *attributes) Strings(key string, values []string) {
	*a = append(*a, attribute.StringSlice(key, values))
}

func (a *attributes) Any(key string, value interface{}) {
	switch value.(type) {
	case []int64:
		v := value.([]int64)
		*a = append(*a, attribute.Int64Slice(key, v))
	case []float64:
		v := value.([]float64)
		*a = append(*a, attribute.Float64Slice(key, v))
	}
}
