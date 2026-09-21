//go:build !noedn

package ordo_test

import (
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestSaveAndLoadEDN(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.edn")

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
