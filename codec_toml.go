//go:build !notoml

package ordo

import (
	"github.com/pelletier/go-toml/v2"
)

type tomlFormat struct{}

func (tomlFormat) Name() string                    { return "toml" }
func (tomlFormat) Extensions() []string            { return []string{".toml"} }
func (tomlFormat) Priority() int                   { return 30 }
func (tomlFormat) Unmarshal(b []byte, v any) error { return toml.Unmarshal(b, v) }
func (tomlFormat) Marshal(v any) ([]byte, error)   { return toml.Marshal(v) }

func init() { RegisterFormat(tomlFormat{}) }
