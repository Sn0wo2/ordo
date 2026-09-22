//go:build !noxml

package ordo

import "encoding/xml"

var XML = NewSimpleFormat("xml", []string{".xml"}, 80, xml.Unmarshal, xml.Marshal)
