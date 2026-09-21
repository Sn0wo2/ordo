package flat

import (
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

func Decode(src, dst any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		MatchName: func(mapKey, fieldName string) bool {
			return strings.EqualFold(mapKey, fieldName)
		},
		Result: dst,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(src)
}
