package ordo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sn0wo2/ordo/internal/utils"
)

func Load[T any](path string, formats ...Format) (*T, error) {
	if len(formats) == 0 {
		return nil, errors.New("no config formats provided")
	}

	path = utils.ResolvePath(path, formats)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	if f, ok := utils.FormatForExtension(formats, filepath.Ext(path)); ok {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	var errs []error

	for _, f := range utils.SortedFormats(formats) {
		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err == nil {
			return cfg, nil
		} else {
			errs = append(errs, fmt.Errorf("%s: %w", f.Name(), err))
		}
	}

	return nil, errors.Join(errs...)
}

func Save(v any, path string, formats ...Format) error {
	if len(formats) == 0 {
		return errors.New("no config formats provided")
	}

	f, ok := utils.FormatForExtension(formats, filepath.Ext(path))
	if !ok {
		f = utils.SortedFormats(formats)[0]
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
