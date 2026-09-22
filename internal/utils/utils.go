package utils

import (
	"cmp"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

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

func SortedFormats(formats []Format) []Format {
	ordered := slices.Clone(formats)

	slices.SortFunc(ordered, func(a, b Format) int {
		if c := cmp.Compare(a.Priority(), b.Priority()); c != 0 {
			return c
		}

		return strings.Compare(a.Name(), b.Name())
	})

	return ordered
}

func FormatForExtension(formats []Format, ext string) (Format, bool) {
	ext = strings.ToLower(ext)

	for _, f := range formats {
		for _, e := range f.Extensions() {
			if strings.ToLower(e) == ext {
				return f, true
			}
		}
	}

	return nil, false
}

func ResolvePath(path string, formats []Format) string {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return path
	}

	base := strings.TrimSuffix(path, filepath.Ext(path))

	seen := make(map[string]bool)
	for _, f := range SortedFormats(formats) {
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
