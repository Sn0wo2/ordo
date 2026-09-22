package ordo

import (
	"encoding/json"
)

type jsonFormat struct{}

func (jsonFormat) Name() string                    { return "json" }
func (jsonFormat) Extensions() []string            { return []string{".json"} }
func (jsonFormat) Priority() int                   { return 20 }
func (jsonFormat) Unmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
func (jsonFormat) Marshal(v any) ([]byte, error)   { return json.Marshal(v) }

func init() { RegisterFormat(jsonFormat{}) }
