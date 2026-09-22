package ordo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sn0wo2/ordo/internal/utils"
)

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

func Load[T any](path string) (*T, error) {
	path = utils.ResolvePath(path)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	if f, ok := utils.ForExtension(filepath.Ext(path)); ok {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	var errs []error

	for _, f := range utils.All() {
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
	f, ok := utils.ForExtension(filepath.Ext(path))
	if !ok {
		available := utils.All()
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
