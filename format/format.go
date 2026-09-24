package format

type Format[T any] interface {
	Marshal(v *T) ([]byte, error)

	Name() string

	Extensions() []string

	Priority() int

	Unmarshal(data []byte, v *T) error
}

type Formatter[T any] struct {
	name      string
	extension []string
	unmarshal func([]byte, any) error
	marshal   func(any) ([]byte, error)
	priority  int
}

func NewFormatter[T any](name string, extension []string, priority int, unmarshal func([]byte, any) error, marshal func(any) ([]byte, error)) Formatter[T] {
	return Formatter[T]{name: name, extension: extension, priority: priority, unmarshal: unmarshal, marshal: marshal}
}

func (f Formatter[T]) Name() string { return f.name }

func (f Formatter[T]) Extensions() []string { return f.extension }

func (f Formatter[T]) Priority() int { return f.priority }

func (f Formatter[T]) Unmarshal(b []byte, v *T) error { return f.unmarshal(b, v) }

func (f Formatter[T]) Marshal(v *T) ([]byte, error) { return f.marshal(v) }
