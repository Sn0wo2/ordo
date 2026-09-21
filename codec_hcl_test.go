//go:build !nohcl

package ordo_test

import (
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestLoadHCL(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.hcl"),
		"name = \"cat\"\nport = 3000\n\nnested {\n  enabled = true\n}\n")

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "cat" || cfg.Port != 3000 || !cfg.Nested.Enabled {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestSaveHCL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.hcl")

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
