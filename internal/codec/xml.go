//go:build !noxml

package codec

import (
	"encoding/xml"

	"github.com/Sn0wo2/ordo/internal/utils"
)

var xmlFormat = utils.NewSimpleFormat("xml", []string{".xml"}, 80, xml.Unmarshal, xml.Marshal)

func init() { utils.Register(xmlFormat) }
