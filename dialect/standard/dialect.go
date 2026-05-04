package standard

import "github.com/laacin/inyorm/intr/dialect"

type DialectStandard struct{}

func LoadStandardDialect() dialect.Dialect {
	return &DialectStandard{}
}
