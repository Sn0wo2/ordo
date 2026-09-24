package ordo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sn0wo2/ordo/format"
	"github.com/Sn0wo2/ordo/internal/utils"
)

type Loader[T any] struct {
	Formats []format.Format

	Default *T

	Merge func(defaultCfg *T, f format.Format, data []byte) (*T, error)

	Defaults func(cfg *T)

	Validate func(cfg *T) error

	OnLoaded func(path string)
}

func (l *Loader[T]) Load(path string) (*T, string, error) {
	if len(l.Formats) == 0 {
		return nil, path, errors.New("no config formats provided")
	}

	path = utils.ResolvePath(path, l.Formats)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, path, fmt.Errorf("failed to load config file %s: %w", path, err)
	}

	cfg, err := func() (*T, error) {
		f, ok := utils.FormatForExtension(l.Formats, filepath.Ext(path))

		if !ok && l.Merge == nil {
			var err error
			if _, f, err = utils.Decode[T](data, path, l.Formats); err != nil {
				return nil, err
			}
		}

		base := new(T)
		if l.Default != nil {
			var err error
			if base, err = utils.DeepCopy(l.Default); err != nil {
				return nil, fmt.Errorf("failed to copy default config: %w", err)
			}
		}

		if l.Merge != nil {
			return l.Merge(base, f, data)
		}

		return base, f.Unmarshal(data, base)
	}()
	if err != nil {
		return nil, path, err
	}

	if l.Defaults != nil {
		l.Defaults(cfg)
	}

	if l.Validate != nil {
		if err := l.Validate(cfg); err != nil {
			return cfg, path, fmt.Errorf("validation failed for config file %s: %w", path, err)
		}
	}

	if l.OnLoaded != nil {
		l.OnLoaded(path)
	}

	return cfg, path, nil
}

func (l *Loader[T]) Save(cfg *T, path string) error {
	return utils.Save(cfg, path, l.Formats...)
}

func Load[T any](path string, formats ...format.Format) (*T, error) {
	if len(formats) == 0 {
		return nil, errors.New("no config formats provided")
	}

	path = utils.ResolvePath(path, formats)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	cfg, _, err := utils.Decode[T](data, path, formats)
	return cfg, err
}
