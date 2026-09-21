// Package ordo is a generic configuration loader with pluggable formats.
//
// JSON is always available (stdlib); YAML and TOML are provided behind
// opt-out build tags: use "-tags noyaml" or "-tags notoml" to exclude a
// format (and its third-party dependency) from the build.
package ordo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Format decodes and encodes a configuration serialization format.
type Format interface {
	// Name returns the canonical name of the format, e.g. "yaml".
	Name() string

	// Extensions returns the file extensions mapped to this format,
	// lowercase and with a leading dot, e.g. [".yml", ".yaml"].
	Extensions() []string

	// Priority orders formats when probing a file whose extension is not
	// registered: lower values are tried first.
	Priority() int

	// Unmarshal decodes data into v.
	Unmarshal(data []byte, v any) error

	// Marshal encodes v.
	Marshal(v any) ([]byte, error)
}

var (
	formatsMu sync.RWMutex
	formats   = map[string]Format{} // lowercase extension (with dot) -> format
)

// RegisterFormat registers f for each of its extensions, overwriting any
// format previously registered for the same extension. It is safe for
// concurrent use.
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

// registeredFormats returns the distinct registered formats ordered by
// Priority, then Name.
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

	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Priority() != ordered[j].Priority() {
			return ordered[i].Priority() < ordered[j].Priority()
		}

		return ordered[i].Name() < ordered[j].Name()
	})

	return ordered
}

// Load reads the file at path and decodes it into T. The format is chosen
// by file extension; an unregistered extension is probed against all
// registered formats in priority order until one succeeds.
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

// Save marshals v and writes it to path, creating parent directories as
// needed (dir mode 0o750, file mode 0o600). The format is chosen by file
// extension; an unregistered extension falls back to the highest-priority
// registered format (yaml, when available).
func Save(v any, path string) error {
	f, ok := formatForExtension(filepath.Ext(path))
	if !ok {
		available := registeredFormats()
		if len(available) == 0 {
			return errors.New("no config format registered")
		}

		f = available[0]
	}

	data, err := f.Marshal(v)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(path), data, 0o600)
}
