//go:build !noini

package ordo

import (
	"bytes"
	"reflect"

	"gopkg.in/ini.v1"

	"github.com/Sn0wo2/ordo/internal/flat"
)

type iniFormat struct{}

func (iniFormat) Name() string         { return "ini" }
func (iniFormat) Extensions() []string { return []string{".ini"} }
func (iniFormat) Priority() int        { return 40 }

func (iniFormat) Unmarshal(b []byte, v any) error {
	f, err := ini.Load(b)
	if err != nil {
		return err
	}

	m := make(map[string]any)

	for _, key := range f.Section(ini.DefaultSection).Keys() {
		m[key.Name()] = key.String()
	}

	for _, name := range f.SectionStrings() {
		if name == ini.DefaultSection {
			continue
		}

		section := make(map[string]any)
		for _, key := range f.Section(name).Keys() {
			section[key.Name()] = key.String()
		}

		m[name] = section
	}

	return flat.Decode(m, v)
}

func (iniFormat) Marshal(v any) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		p := reflect.New(rv.Type())
		p.Elem().Set(rv)
		rv = p
	}

	file := ini.Empty()
	if err := file.ReflectFrom(rv.Interface()); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := file.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func init() { RegisterFormat(iniFormat{}) }
