package format

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

type Format interface {
	Marshaler

	Name() string

	Extensions() []string

	Priority() int

	Unmarshal(data []byte, v any) error
}

type Formatter struct {
	name      string
	extension []string
	unmarshal func([]byte, any) error
	marshal   func(any) ([]byte, error)
	priority  int
}

func NewFormatter(name string, extension []string, priority int, unmarshal func([]byte, any) error, marshal func(any) ([]byte, error)) Formatter {
	return Formatter{name: name, extension: extension, priority: priority, unmarshal: unmarshal, marshal: marshal}
}

func (f Formatter) Name() string { return f.name }

func (f Formatter) Extensions() []string { return f.extension }

func (f Formatter) Priority() int { return f.priority }

func (f Formatter) Unmarshal(b []byte, v any) error { return f.unmarshal(b, v) }

func (f Formatter) Marshal(v any) ([]byte, error) { return f.marshal(v) }
