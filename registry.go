package ordo

import "github.com/Sn0wo2/ordo/internal/utils"

type Format = utils.Format

func NewSimpleFormat(name string, extension []string, priority int, unmarshal func([]byte, any) error, marshal func(any) ([]byte, error)) Format {
	return utils.NewSimpleFormat(name, extension, priority, unmarshal, marshal)
}

func RegisterFormat(f Format) { utils.Register(f) }
