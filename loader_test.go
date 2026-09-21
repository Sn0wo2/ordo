package ordo_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Sn0wo2/ordo"
)

type testNested struct {
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`
}

type testConfig struct {
	Name   string     `json:"name"   yaml:"name"   toml:"name"`
	Port   int        `json:"port"   yaml:"port"   toml:"port"`
	Nested testNested `json:"nested" yaml:"nested" toml:"nested"`
}

func writeTestFile(t *testing.T, path, content string) string {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	return path
}

func TestLoadJSON(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.json"),
		`{"name":"cat","port":3000,"nested":{"enabled":true}}`)

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "cat" || cfg.Port != 3000 || !cfg.Nested.Enabled {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadUnknownExtensionFallback(t *testing.T) {
	// Unknown extension: registered formats are probed in priority order
	// until one succeeds (JSON is always registered).
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.conf"),
		`{"name":"dog","port":8080,"nested":{"enabled":false}}`)

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "dog" || cfg.Port != 8080 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadInvalidContent(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.json"), `{invalid`)

	if _, err := ordo.Load[testConfig](path); err == nil {
		t.Fatal("expected error for invalid content")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := ordo.Load[testConfig](filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got: %v", err)
	}
}

func TestSaveCreatesDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "config.json")

	if err := ordo.Save(testConfig{Name: "saved", Port: 1}, path); err != nil {
		t.Fatalf("save: %v", err)
	}

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load back: %v", err)
	}

	if cfg.Name != "saved" || cfg.Port != 1 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoaderHooks(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.json"), `{"name":"raw"}`)

	var (
		defaultsCalled bool
		loadedPath     string
	)

	loader := &ordo.Loader[testConfig]{
		Defaults: func(cfg *testConfig) {
			defaultsCalled = true

			if cfg.Port == 0 {
				cfg.Port = 9999
			}
		},
		Validate: func(cfg *testConfig) error {
			if cfg.Name == "" {
				return errors.New("name is required")
			}

			return nil
		},
		OnLoaded: func(p string) { loadedPath = p },
	}

	cfg, gotPath, err := loader.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if !defaultsCalled {
		t.Fatal("defaults hook was not called")
	}

	if cfg.Port != 9999 {
		t.Fatalf("defaults were not applied: %+v", cfg)
	}

	if gotPath != path || loadedPath != path {
		t.Fatalf("unexpected paths: %q / %q", gotPath, loadedPath)
	}
}

func TestLoaderValidateError(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.json"), `{"port":1}`)

	loader := &ordo.Loader[testConfig]{
		Validate: func(cfg *testConfig) error {
			return errors.New("name is required")
		},
		OnLoaded: func(string) { t.Fatal("onLoaded must not run when validation fails") },
	}

	cfg, _, err := loader.Load(path)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "validation failed for config file") {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("config must be returned even when validation fails")
	}
}

func TestLoaderLoadError(t *testing.T) {
	loader := &ordo.Loader[testConfig]{}

	_, _, err := loader.Load(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "failed to load config file") {
		t.Fatalf("unexpected error: %v", err)
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("os.ErrNotExist must stay detectable, got: %v", err)
	}
}

func TestNilLoaderHooksAreSkipped(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.json"), `{}`)

	var loader ordo.Loader[testConfig]

	if _, _, err := loader.Load(path); err != nil {
		t.Fatalf("load: %v", err)
	}
}
