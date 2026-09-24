package format

import "encoding/xml"

var XML = NewFormatter("xml", []string{".xml"}, 80, xml.Unmarshal, xml.Marshal)
