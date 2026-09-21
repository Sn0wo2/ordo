//go:build !noenv

package ordo

import (
	"github.com/joho/godotenv"
)

type envFormat struct{}

func (envFormat) Name() string         { return "env" }
func (envFormat) Extensions() []string { return []string{".env"} }
func (envFormat) Priority() int        { return 50 }

// Unmarshal decodes a flat KEY=value document into a struct. Keys are matched
// against json, then yaml, then Go field names, case-insensitively; values
// are assigned only when the field type accepts them. The env format is
// read-only: it does not implement Marshaler.
func (envFormat) Unmarshal(b []byte, v any) error {
	m, err := godotenv.Unmarshal(string(b))
	if err != nil {
		return err
	}

	flat := make(map[string]any, len(m))
	for k, val := range m {
		flat[k] = val
	}

	return decodeFlatMap(flat, v)
}

func init() { RegisterFormat(envFormat{}) }
