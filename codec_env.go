//go:build !noenv

package ordo

import (
	"github.com/joho/godotenv"
)

type envFormat struct{}

func (envFormat) Name() string         { return "env" }
func (envFormat) Extensions() []string { return []string{".env"} }
func (envFormat) Priority() int        { return 50 }

func (envFormat) Unmarshal(b []byte, v any) error {
	m, err := godotenv.Unmarshal(string(b))
	if err != nil {
		return err
	}

	return decodeFlat(m, v)
}

func init() { RegisterFormat(envFormat{}) }
