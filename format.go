package ordo

import "github.com/Sn0wo2/ordo/internal/utils"

type (
	Marshaler = utils.Marshaler

	Format = utils.Format

	SimpleFormat = utils.SimpleFormat
)

func NewSimpleFormat(name string, extension []string, priority int, unmarshal func([]byte, any) error, marshal func(any) ([]byte, error)) SimpleFormat {
	return utils.NewSimpleFormat(name, extension, priority, unmarshal, marshal)
}
