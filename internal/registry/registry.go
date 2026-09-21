package registry

import (
	"cmp"
	"slices"
	"strings"
	"sync"
)

type Format interface {
	Name() string

	Extensions() []string

	Priority() int

	Unmarshal(data []byte, v any) error
}

var (
	mu          sync.RWMutex
	byExtension = map[string]Format{}
)

func Register(f Format) {
	mu.Lock()
	defer mu.Unlock()

	for _, ext := range f.Extensions() {
		byExtension[strings.ToLower(ext)] = f
	}
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
