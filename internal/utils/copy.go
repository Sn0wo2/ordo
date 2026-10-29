package utils

import (
	"fmt"
	"reflect"
)

func DeepCopy[T any](src *T) (*T, error) {
	if c, ok := any(src).(interface{ DeepCopy() T }); ok {
		dst := c.DeepCopy()

		return &dst, nil
	}

	var copyValue func(reflect.Value) (reflect.Value, error)

	copyValue = func(src reflect.Value) (reflect.Value, error) {
		needsDeepCopy := func(k reflect.Kind) bool {
			return k == reflect.Pointer || k == reflect.Interface || k == reflect.Slice ||
				k == reflect.Map || k == reflect.Array || k == reflect.Struct ||
				k == reflect.Func || k == reflect.Chan || k == reflect.UnsafePointer
		}

		switch src.Kind() {
		case reflect.Chan, reflect.Func, reflect.UnsafePointer:
			return reflect.Value{}, fmt.Errorf("unsupported kind %s", src.Kind())
		case reflect.Pointer:
			if src.IsNil() {
				return src, nil
			}

			dst := reflect.New(src.Type().Elem())

			elem, err := copyValue(src.Elem())
			if err != nil {
				return reflect.Value{}, err
			}

			dst.Elem().Set(elem)

			return dst, nil
		case reflect.Interface:
			if src.IsNil() {
				return src, nil
			}

			elem, err := copyValue(src.Elem())
			if err != nil {
				return reflect.Value{}, err
			}

			dst := reflect.New(src.Type()).Elem()
			dst.Set(elem)

			return dst, nil
		case reflect.Slice:
			if src.IsNil() {
				return src, nil
			}

			dst := reflect.MakeSlice(src.Type(), src.Len(), src.Len())
			for i := range src.Len() {
				elem, err := copyValue(src.Index(i))
				if err != nil {
					return reflect.Value{}, err
				}

				dst.Index(i).Set(elem)
			}

			return dst, nil
		case reflect.Map:
			if src.IsNil() {
				return src, nil
			}

			dst := reflect.MakeMapWithSize(src.Type(), src.Len())

			iter := src.MapRange()
			for iter.Next() {
				key, err := copyValue(iter.Key())
				if err != nil {
					return reflect.Value{}, err
				}

				val, err := copyValue(iter.Value())
				if err != nil {
					return reflect.Value{}, err
				}

				dst.SetMapIndex(key, val)
			}

			return dst, nil
		case reflect.Array:
			dst := reflect.New(src.Type()).Elem()
			for i := range src.Len() {
				elem, err := copyValue(src.Index(i))
				if err != nil {
					return reflect.Value{}, err
				}

				dst.Index(i).Set(elem)
			}

			return dst, nil
		case reflect.Struct:
			dst := reflect.New(src.Type()).Elem()
			dst.Set(src)

			for i := range src.NumField() {
				if !dst.Field(i).CanSet() || !needsDeepCopy(src.Field(i).Kind()) {
					continue
				}

				elem, err := copyValue(src.Field(i))
				if err != nil {
					return reflect.Value{}, err
				}

				dst.Field(i).Set(elem)
			}

			return dst, nil
		case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64,
			reflect.Complex64, reflect.Complex128, reflect.String:
			return src, nil
		}

		return src, nil
	}

	v, err := copyValue(reflect.ValueOf(src))
	if err != nil {
		return nil, err
	}

	dst, _ := v.Interface().(*T)

	return dst, nil
}
