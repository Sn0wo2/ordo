package utils

import (
	"cmp"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

type DecodeOptions struct {
	StrictTypes bool
}

type Format interface {
	Marshaler

	Name() string

	Extensions() []string

	Priority() int

	Unmarshal(data []byte, v any) error

	WithOption(opt DecodeOptions) Format
}

type SimpleFormat struct {
	name      string
	extension []string
	priority  int
	unmarshal func([]byte, any) error
	marshal   func(any) ([]byte, error)
}

func NewSimpleFormat(name string, extension []string, priority int, unmarshal func([]byte, any) error, marshal func(any) ([]byte, error)) SimpleFormat {
	return SimpleFormat{name: name, extension: extension, priority: priority, unmarshal: unmarshal, marshal: marshal}
}

func (f SimpleFormat) Name() string                    { return f.name }
func (f SimpleFormat) Extensions() []string            { return f.extension }
func (f SimpleFormat) Priority() int                   { return f.priority }
func (f SimpleFormat) Unmarshal(b []byte, v any) error { return f.unmarshal(b, v) }
func (f SimpleFormat) Marshal(v any) ([]byte, error)   { return f.marshal(v) }

func (f SimpleFormat) WithOption(DecodeOptions) Format { return f }

var (
	mu          sync.RWMutex
	byExtension = map[string]Format{}
	byName      = map[string]Format{}
)

func Register(f Format) {
	mu.Lock()
	defer mu.Unlock()

	byName[f.Name()] = f

	for _, ext := range f.Extensions() {
		byExtension[strings.ToLower(ext)] = f
	}
}

func ByName(name string) (Format, bool) {
	mu.RLock()
	defer mu.RUnlock()

	f, ok := byName[name]

	return f, ok
}

func ForExtension(ext string) (Format, bool) {
	mu.RLock()
	defer mu.RUnlock()

	f, ok := byExtension[strings.ToLower(ext)]

	return f, ok
}

func All() []Format {
	mu.RLock()
	defer mu.RUnlock()

	seen := make(map[string]Format, len(byExtension))
	for _, f := range byExtension {
		seen[f.Name()] = f
	}

	ordered := make([]Format, 0, len(seen))
	for _, f := range seen {
		ordered = append(ordered, f)
	}

	slices.SortFunc(ordered, func(a, b Format) int {
		if c := cmp.Compare(a.Priority(), b.Priority()); c != 0 {
			return c
		}

		return strings.Compare(a.Name(), b.Name())
	})

	return ordered
}

func ResolvePath(path string) string {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return path
	}

	base := strings.TrimSuffix(path, filepath.Ext(path))

	seen := make(map[string]bool)
	for _, f := range All() {
		for _, ext := range f.Extensions() {
			candidate := base + ext
			if seen[candidate] {
				continue
			}

			seen[candidate] = true

			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}

	return path
}
