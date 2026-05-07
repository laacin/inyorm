package standard

import "github.com/laacin/inyorm/internal/entity"

type DialectStandard struct{}

func DialectDefault() entity.Dialect {
	return &DialectStandard{}
}
