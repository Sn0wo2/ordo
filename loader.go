package ordo

import "fmt"

type Loader[T any] struct {
	Defaults func(cfg *T)

	// Validate, when set, checks the decoded (and defaulted) config.
	Validate func(cfg *T) error

	// OnLoaded, when set, is called with the path after a successful load.
	OnLoaded func(path string)
}

func (l *Loader[T]) Load(path string) (*T, string, error) {
	cfg, err := Load[T](path)
	if err != nil {
		return nil, path, fmt.Errorf("failed to load config file %s: %w", path, err)
	}

	if l != nil {
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
	}

	return cfg, path, nil
}
