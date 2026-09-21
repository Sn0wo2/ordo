//go:build !noenv

package ordo_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestLoadENV(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.env"),
		"NAME=cat\nPORT=3000\n")

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "cat" || cfg.Port != 3000 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadENVIgnoresUnknownKeys(t *testing.T) {
	path := writeTestFile(t, filepath.Join(t.TempDir(), "config.env"),
		"NAME=cat\nUNKNOWN_FIELD=nope\n")

	cfg, err := ordo.Load[testConfig](path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.Name != "cat" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestSaveENVUnsupported(t *testing.T) {
	// env is flat and read-only: Save must refuse it.
	err := ordo.Save(testConfig{Name: "dog"}, filepath.Join(t.TempDir(), "config.env"))
	if err == nil {
		t.Fatal("expected error when saving env format")
	}

	if !strings.Contains(err.Error(), "does not support saving") {
		t.Fatalf("unexpected error: %v", err)
	}
}
