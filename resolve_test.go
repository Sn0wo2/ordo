package ordo_test

import (
	"path/filepath"
	"testing"

	"github.com/Sn0wo2/ordo"
)

func TestResolvePathExplicitWins(t *testing.T) {
	// An explicit path is returned unchanged, even when it does not exist.
	got := ordo.ResolvePath("explicit.yml", "preferred.json", "default.yml")
	if got != "explicit.yml" {
		t.Fatalf("expected explicit path, got %q", got)
	}
}

func TestResolvePathDefaultWhenNoPreferred(t *testing.T) {
	got := ordo.ResolvePath("", "", "default.yml")
	if got != "default.yml" {
		t.Fatalf("expected default path, got %q", got)
	}
}

func TestResolvePathPreferredExists(t *testing.T) {
	dir := t.TempDir()
	preferred := filepath.Join(dir, "config.json")
	writeTestFile(t, preferred, `{}`)

	got := ordo.ResolvePath("", preferred, "default.yml")
	if got != preferred {
		t.Fatalf("expected %q, got %q", preferred, got)
	}
}

func TestResolvePathProbesRegisteredExtensions(t *testing.T) {
	dir := t.TempDir()
	found := filepath.Join(dir, "config.json")
	writeTestFile(t, found, `{}`)

	// The preferred path has no extension; the same base name is probed
	// with every registered format's extension.
	got := ordo.ResolvePath("", filepath.Join(dir, "config"), "default.yml")
	if got != found {
		t.Fatalf("expected %q, got %q", found, got)
	}
}

func TestResolvePathProbesSiblingExtension(t *testing.T) {
	dir := t.TempDir()
	found := filepath.Join(dir, "config.json")
	writeTestFile(t, found, `{}`)

	// The preferred path names a file whose extension does not exist on
	// disk; a sibling with another registered extension is found.
	got := ordo.ResolvePath("", filepath.Join(dir, "config.toml"), "default.yml")
	if got != found {
		t.Fatalf("expected %q, got %q", found, got)
	}
}

func TestResolvePathFallsBackToPreferred(t *testing.T) {
	dir := t.TempDir()
	preferred := filepath.Join(dir, "missing")

	got := ordo.ResolvePath("", preferred, "default.yml")
	if got != preferred {
		t.Fatalf("expected preferred path, got %q", got)
	}
}
