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
			switch k {
			case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map,
				reflect.Array, reflect.Struct, reflect.Func, reflect.Chan, reflect.UnsafePointer:
				return true
			default:
				return false
			}
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
			for i := 0; i < src.Len(); i++ {
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
			for i := 0; i < src.Len(); i++ {
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

			for i := 0; i < src.NumField(); i++ {
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
		default:
			return src, nil
		}
	}

	v, err := copyValue(reflect.ValueOf(src))
	if err != nil {
		return nil, err
	}

	dst, _ := v.Interface().(*T)
	return dst, nil
}
