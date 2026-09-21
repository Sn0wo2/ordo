package ordo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Sn0wo2/ordo/internal/registry"
)

func ResolvePath(explicitPath, preferredPath, defaultPath string) string {
	if explicitPath != "" {
		return explicitPath
	}

	if preferredPath == "" {
		return defaultPath
	}

	if info, err := os.Stat(preferredPath); err == nil && !info.IsDir() {
		return preferredPath
	}

	base := strings.TrimSuffix(preferredPath, filepath.Ext(preferredPath))

	seen := make(map[string]bool)
	for _, f := range registry.All() {
		for _, ext := range f.Extensions() {
			candidate := base + ext
			if seen[candidate] {
				continue
			}

			seen[candidate] = true

			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}

	return preferredPath
}
