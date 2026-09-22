package ordo

import "github.com/Sn0wo2/ordo/internal/utils"

type Format = utils.Format

func RegisterFormat(f Format) { utils.Register(f) }
