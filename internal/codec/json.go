package codec

import (
	"encoding/json"

	"github.com/Sn0wo2/ordo/internal/utils"
)

var jsonFormat = utils.NewSimpleFormat("json", []string{".json"}, 20, json.Unmarshal, json.Marshal)

func init() { utils.Register(jsonFormat) }
