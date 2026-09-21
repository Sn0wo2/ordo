package ordo

import "fmt"

// Loader orchestrates loading a configuration file: decode, apply defaults,
// validate and notify. Every hook is optional; a nil hook is skipped. ordo
// knows nothing about the concrete config type beyond the type parameter.
type Loader[T any] struct {
	// Defaults, when set, fills unset fields of the decoded config in place,
	// before Validate runs.
	Defaults func(cfg *T)

	// Validate, when set, checks the decoded (and defaulted) config.
	Validate func(cfg *T) error

	// OnLoaded, when set, is called with the path after a successful load.
	OnLoaded func(path string)
}

// Load reads and decodes the config file at path, then runs the Defaults,
// Validate and OnLoaded hooks. The path is returned alongside the config so
// callers can persist changes back to the same file.
//
// Errors wrap the underlying cause: os.ErrNotExist stays detectable via
// errors.Is, and validation failures are wrapped with the path included.
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
