package ordo_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestLoadProbesGivenExtensions(t *testing.T) {
	dir := t.TempDir()
	found := filepath.Join(dir, "config.json")
	writeTestFile(t, found, `{}`)

	if _, err := ordo.Load[testConfig](filepath.Join(dir, "config"), testFormats...); err != nil {
		t.Fatalf("load: %v", err)
	}
}

func TestLoadProbesSiblingExtension(t *testing.T) {
	dir := t.TempDir()
	found := filepath.Join(dir, "config.json")
	writeTestFile(t, found, `{}`)

	if _, err := ordo.Load[testConfig](filepath.Join(dir, "config.toml"), testFormats...); err != nil {
		t.Fatalf("load: %v", err)
	}
}

func TestLoadMissingPathFallsBackToError(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing")

	if _, err := ordo.Load[testConfig](missing, testFormats...); err == nil {
		t.Fatal("expected error for missing file")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got: %v", err)
	}
}
