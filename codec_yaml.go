//go:build !noyaml

package ordo

import (
	"gopkg.in/yaml.v3"
)

type yamlFormat struct{}

func (yamlFormat) Name() string                    { return "yaml" }
func (yamlFormat) Extensions() []string            { return []string{".yaml", ".yml"} }
func (yamlFormat) Priority() int                   { return 10 }
func (yamlFormat) Unmarshal(b []byte, v any) error { return yaml.Unmarshal(b, v) }
func (yamlFormat) Marshal(v any) ([]byte, error)   { return yaml.Marshal(v) }

func init() { RegisterFormat(yamlFormat{}) }
