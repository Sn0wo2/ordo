//go:build !noedn

package ordo

import (
	"olympos.io/encoding/edn"
)

type ednFormat struct{}

func (ednFormat) Name() string         { return "edn" }
func (ednFormat) Extensions() []string { return []string{".edn"} }
func (ednFormat) Priority() int        { return 70 }

func (ednFormat) Unmarshal(b []byte, v any) error { return edn.Unmarshal(b, v) }
func (ednFormat) Marshal(v any) ([]byte, error)   { return edn.Marshal(v) }

func init() { RegisterFormat(ednFormat{}) }
