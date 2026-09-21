package ordo

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// decodeFlatMap assigns a map with (mostly) scalar leaves into a struct
// value. Keys are matched against json, then yaml, then Go field names,
// case-insensitively. String values are coerced to the field type; nested
// map[string]any values are assigned to struct (or pointer-to-struct)
// fields. Fields without a matching key are left untouched.
func decodeFlatMap(m map[string]any, v any) error {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return errors.New("nil destination")
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return errors.New("destination must be a struct")
	}

	return flatAssignStruct(rv, m)
}

func flatFieldName(field reflect.StructField) string {
	for _, tagKey := range []string{"json", "yaml"} {
		tag, ok := field.Tag.Lookup(tagKey)
		if !ok {
			continue
		}

		if name := strings.Split(tag, ",")[0]; name != "" && name != "-" {
			return name
		}
	}

	return field.Name
}

func flatLookup(m map[string]any, key string) (any, bool) {
	if val, ok := m[key]; ok {
		return val, true
	}

	upper := strings.ToUpper(key)
	for k, val := range m {
		if strings.ToUpper(k) == upper {
			return val, true
		}
	}

	return nil, false
}

func flatAssignStruct(rv reflect.Value, m map[string]any) error {
	t := rv.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		key := flatFieldName(field)

		val, ok := flatLookup(m, key)
		if !ok {
			continue
		}

		if err := flatAssignValue(rv.Field(i), val); err != nil {
			return fmt.Errorf("key %q: %w", key, err)
		}
	}

	return nil
}

func flatAssignValue(fv reflect.Value, val any) error {
	if fv.Kind() == reflect.Pointer {
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}

		return flatAssignValue(fv.Elem(), val)
	}

	switch typed := val.(type) {
	case map[string]any:
		if fv.Kind() != reflect.Struct {
			return fmt.Errorf("cannot assign map to %s field", fv.Kind())
		}

		return flatAssignStruct(fv, typed)
	case string:
		return flatAssignString(fv, typed)
	case bool:
		if fv.Kind() != reflect.Bool {
			return fmt.Errorf("cannot assign bool to %s field", fv.Kind())
		}

		fv.SetBool(typed)
	case int64:
		return flatAssignInt(fv, typed, nil)
	case float64:
		switch fv.Kind() {
		case reflect.Float32, reflect.Float64:
			fv.SetFloat(typed)
		default:
			return flatAssignInt(fv, 0, &typed)
		}
	default:
		return fmt.Errorf("unsupported value type %T", val)
	}

	return nil
}

func flatAssignString(fv reflect.Value, val string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(val)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}

		fv.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(val, 10, fv.Type().Bits())
		if err != nil {
			return err
		}

		fv.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(val, 10, fv.Type().Bits())
		if err != nil {
			return err
		}

		fv.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(val, fv.Type().Bits())
		if err != nil {
			return err
		}

		fv.SetFloat(parsed)
	default:
		return fmt.Errorf("unsupported field kind %s", fv.Kind())
	}

	return nil
}

func flatAssignInt(fv reflect.Value, i int64, f *float64) error {
	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if f != nil {
			i = int64(*f)
		}

		fv.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if f != nil {
			i = int64(*f)
		}

		if i < 0 {
			return fmt.Errorf("negative value %d into unsigned field", i)
		}

		fv.SetUint(uint64(i))
	default:
		return fmt.Errorf("cannot assign number to %s field", fv.Kind())
	}

	return nil
}
