package ordo

func Unmarshal[T any](f Format, data []byte) (*T, error) {
	v := new(T)
	if err := f.Unmarshal(data, v); err != nil {
		return nil, err
	}

	return v, nil
}
