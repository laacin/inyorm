package builder

import (
	"github.com/laacin/inyorm/internal/builder/mapper"
	"github.com/laacin/inyorm/internal/builder/params"
	"github.com/laacin/inyorm/internal/core"
)

func New() *core.Builder {
	return &core.Builder{
		Mapper: &mapper.Mapper{},
		Param:  &params.ParamStore{},
	}
}
