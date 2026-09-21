package ordo

import (
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

// decodeFlat copies a flat map of config values (env or ini style, where all
// values start out as strings) into a struct pointer. Field values are
// converted weakly ("3000" -> int), and keys match field names or their json
// tag case-insensitively so UPPER_CASE env keys find their fields.
func decodeFlat(src, dst any) error {
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
