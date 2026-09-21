//go:build !noini

package ordo

import (
	"bytes"
	"reflect"

	"gopkg.in/ini.v1"
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

	return decodeFlatMap(iniToMap(f), v)
}

// iniToMap flattens an ini file into a map: keys of the default section end
// up at the top level, every other section becomes a nested map. All values
// are strings; type assignment is handled by decodeFlatMap.
func iniToMap(f *ini.File) map[string]any {
	out := make(map[string]any)

	for _, key := range f.Section(ini.DefaultSection).Keys() {
		out[key.Name()] = key.String()
	}

	for _, name := range f.SectionStrings() {
		if name == ini.DefaultSection {
			continue
		}

		section := make(map[string]any)
		for _, key := range f.Section(name).Keys() {
			section[key.Name()] = key.String()
		}

		out[name] = section
	}

	return out
}

func (iniFormat) Marshal(v any) ([]byte, error) {
	// ReflectFrom requires a pointer to a struct; wrap values transparently.
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		p := reflect.New(rv.Type())
		p.Elem().Set(rv)
		rv = p
	}

	f := ini.Empty()
	if err := f.ReflectFrom(rv.Interface()); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func init() { RegisterFormat(iniFormat{}) }
