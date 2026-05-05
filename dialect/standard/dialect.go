package standard

import "github.com/laacin/inyorm/internal/entity"

type DialectStandard struct{}

func LoadStandardDialect() entity.Dialect {
	return &DialectStandard{}
}
