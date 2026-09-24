package format

import "encoding/json/v2"

func JSON[T any](opts ...json.Options) Format[T] {
	return NewFormatter[T]("json", []string{".json"}, 10,
		func(b []byte, v any) error { return json.Unmarshal(b, v, opts...) },
		func(v any) ([]byte, error) { return json.Marshal(v, opts...) },
	)
}
