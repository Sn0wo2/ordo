package format

import "encoding/json"

var JSON = NewFormatter("json", []string{".json"}, 10, json.Unmarshal, json.Marshal)
