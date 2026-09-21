package ordo

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Format interface {
	Name() string

	// lowercase and with a leading dot, e.g. [".yml", ".yaml"].
	Extensions() []string

	// lower values are tried first.
	Priority() int

	Unmarshal(data []byte, v any) error
}

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

var (
	formatsMu sync.RWMutex
	formats   = map[string]Format{} // lowercase extension (with dot) -> format
)

func RegisterFormat(f Format) {
	formatsMu.Lock()
	defer formatsMu.Unlock()

	for _, ext := range f.Extensions() {
		formats[strings.ToLower(ext)] = f
	}
}

func formatForExtension(ext string) (Format, bool) {
	formatsMu.RLock()
	defer formatsMu.RUnlock()

	f, ok := formats[strings.ToLower(ext)]

	return f, ok
}

func registeredFormats() []Format {
	formatsMu.RLock()
	defer formatsMu.RUnlock()

	seen := make(map[string]Format, len(formats))
	for _, f := range formats {
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

func Load[T any](path string) (*T, error) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	if f, ok := formatForExtension(filepath.Ext(path)); ok {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	var errs []error

	for _, f := range registeredFormats() {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err == nil {
			return cfg, nil
		} else {
			errs = append(errs, fmt.Errorf("%s: %w", f.Name(), err))
		}
	}

	if len(errs) == 0 {
		return nil, errors.New("no config format registered")
	}

	return nil, errors.Join(errs...)
}

func Save(v any, path string) error {
	f, ok := formatForExtension(filepath.Ext(path))
	if !ok {
		available := registeredFormats()
		if len(available) == 0 {
			return errors.New("no config format registered")
		}

		f = available[0]
	}

	marshaler, ok := f.(Marshaler)
	if !ok {
		return fmt.Errorf("format %q does not support saving", f.Name())
	}

	data, err := marshaler.Marshal(v)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(path), data, 0o600)
}
