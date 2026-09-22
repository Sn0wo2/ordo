package ordo

import "encoding/json"

var JSON = NewSimpleFormat("json", []string{".json"}, 20, json.Unmarshal, json.Marshal)
