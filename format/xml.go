package format

import "encoding/xml"

var XML = NewFormatter("xml", []string{".xml"}, 10, xml.Unmarshal, xml.Marshal)
