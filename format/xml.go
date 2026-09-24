package format

import "encoding/xml"

func XML[T any]() Format[T] {
	return NewFormatter[T]("xml", []string{".xml"}, 10, xml.Unmarshal, xml.Marshal)
}
