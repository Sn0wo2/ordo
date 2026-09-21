package ordo

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath determines which config file to load:
//
//  1. A non-empty explicitPath always wins and is returned unchanged.
//  2. Otherwise, when preferredPath is non-empty: the path itself is tried
//     first, then the same base name with each registered format's extension
//     (in format priority order); the first existing file is returned. If
//     nothing exists, preferredPath is returned so loading reports a precise
//     error for the path the caller actually asked for.
//  3. Otherwise defaultPath is returned.
func ResolvePath(explicitPath, preferredPath, defaultPath string) string {
	if explicitPath != "" {
		return explicitPath
	}

	if preferredPath == "" {
		return defaultPath
	}

	if fileExists(preferredPath) {
		return preferredPath
	}

	base := strings.TrimSuffix(preferredPath, filepath.Ext(preferredPath))

	seen := make(map[string]bool)
	for _, f := range registeredFormats() {
		for _, ext := range f.Extensions() {
			candidate := base + ext
			if seen[candidate] {
				continue
			}

			seen[candidate] = true

			if fileExists(candidate) {
				return candidate
			}
		}
	}

	return preferredPath
}

func fileExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}
