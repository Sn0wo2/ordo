//go:build !noini

package ordo_test

import (
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestLoadINI(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.ini"),
		"name = cat\nport = 3000\n\n[nested]\nenabled = true\n")

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "cat" || cfg.Port != 3000 || !cfg.Nested.Enabled {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestSaveINI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.ini")

	if err := ordo.Save(testConfig{Name: "dog", Port: 8080, Nested: testNested{Enabled: true}}, path); err != nil {
		t.Fatalf("save: %v", err)
	}

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load back: %v", err)
	}

	if cfg.Name != "dog" || cfg.Port != 8080 || !cfg.Nested.Enabled {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
