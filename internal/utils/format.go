package utils

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Sn0wo2/ordo/format"
)

func Decode[T any](data []byte, path string, formats []format.Format[T]) (*T, format.Format[T], error) {
	if f, ok := FormatForExtension(formats, filepath.Ext(path)); ok {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err != nil {
			return nil, nil, err
		}

		return cfg, f, nil
	}

	var errs []error

	for _, f := range SortedFormats(formats) {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err == nil {
			return cfg, f, nil
		} else {
			errs = append(errs, fmt.Errorf("%s: %w", f.Name(), err))
		}
	}

	return nil, nil, errors.Join(errs...)
}

func Save[T any](v *T, path string, formats ...format.Format[T]) error {
	if len(formats) == 0 {
		return errors.New("no config formats provided")
	}

	f, ok := FormatForExtension(formats, filepath.Ext(path))
	if !ok {
		f = SortedFormats(formats)[0]
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

func SortedFormats[T any](formats []format.Format[T]) []format.Format[T] {
	ordered := slices.Clone(formats)

	slices.SortFunc(ordered, func(a, b format.Format[T]) int {
		if c := cmp.Compare(a.Priority(), b.Priority()); c != 0 {
			return c
		}

		return strings.Compare(a.Name(), b.Name())
	})

	return ordered
}

func FormatForExtension[T any](formats []format.Format[T], ext string) (format.Format[T], bool) {
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

func ResolvePath[T any](path string, formats []format.Format[T]) string {
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
