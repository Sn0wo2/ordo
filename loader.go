package ordo

import (
	"fmt"

	"github.com/Sn0wo2/ordo/internal/utils"
)

type Loader[T any] struct {
	Defaults func(cfg *T)

	Validate func(cfg *T) error

	OnLoaded func(path string)
}

func (l *Loader[T]) Load(path string) (*T, string, error) {
	path = utils.ResolvePath(path)

	cfg, err := Load[T](path)
	if err != nil {
		return nil, path, fmt.Errorf("failed to load config file %s: %w", path, err)
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
