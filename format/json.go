package format

import "encoding/json"

var JSON = NewFormatter("json", []string{".json"}, 20, json.Unmarshal, json.Marshal)
