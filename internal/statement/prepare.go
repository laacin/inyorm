package statement

import (
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/query"
)

type Prepared struct {
	query  string
	params core.ParamStore
}

func NewPrepared(driver core.Driver, qc *query.Compiler) (*Prepared, error) {
	result, err := qc.Compile()
	if err != nil {
		return nil, err
	}

	return &Prepared{
		query:  result.QueryString,
		params: result.Params,
	}, nil
}
