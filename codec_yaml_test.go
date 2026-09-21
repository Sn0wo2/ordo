//go:build !noyaml

package ordo_test

import (
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestLoadYAML(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.yml"),
		"name: cat\nport: 3000\nnested:\n  enabled: true\n")

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "cat" || cfg.Port != 3000 || !cfg.Nested.Enabled {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestSaveDefaultsToYAML(t *testing.T) {
	// Unknown extension falls back to the highest-priority format (yaml).
	path := filepath.Join(t.TempDir(), "config.cfg")

	if err := ordo.Save(testConfig{Name: "dog", Port: 8080}, path); err != nil {
		t.Fatalf("save: %v", err)
	}

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load back: %v", err)
	}

	if cfg.Name != "dog" || cfg.Port != 8080 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}
