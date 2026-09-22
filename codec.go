package ordo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/Sn0wo2/ordo/internal/codec"
	"github.com/Sn0wo2/ordo/internal/utils"
)

type Marshaler = utils.Marshaler

type loadConfig struct {
	forceFormat string
	decodeOpts  utils.DecodeOptions
}

type LoadOption func(*loadConfig)

func WithFormat(name string) LoadOption {
	return func(c *loadConfig) { c.forceFormat = name }
}

func WithStrictTypes() LoadOption {
	return func(c *loadConfig) { c.decodeOpts.StrictTypes = true }
}

func Load[T any](path string, opts ...LoadOption) (*T, error) {
	var lc loadConfig
	for _, opt := range opts {
		opt(&lc)
	}

	path = utils.ResolvePath(path)

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	if lc.forceFormat != "" {
		f, ok := utils.ByName(lc.forceFormat)
		if !ok {
			return nil, fmt.Errorf("unknown config format %q", lc.forceFormat)
		}

		if lc.decodeOpts != (utils.DecodeOptions{}) {
			f = f.WithOption(lc.decodeOpts)
		}

		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	if f, ok := utils.ForExtension(filepath.Ext(path)); ok {
		if lc.decodeOpts != (utils.DecodeOptions{}) {
			f = f.WithOption(lc.decodeOpts)
		}

		cfg := new(T)
		if err := f.Unmarshal(data, cfg); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	var errs []error

	for _, f := range utils.All() {
		if lc.decodeOpts != (utils.DecodeOptions{}) {
			f = f.WithOption(lc.decodeOpts)
		}

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

	data, err := f.Marshal(v)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(path), data, 0o600)
}
