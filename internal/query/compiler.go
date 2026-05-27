package query

import (
	"fmt"

	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/core/aliases"
	"github.com/laacin/inyorm/internal/core/params"
	"github.com/laacin/inyorm/internal/core/writer"
	"github.com/laacin/inyorm/internal/expr"
)

type Compiler struct {
	tools *Tools
	query Query
}

func NewCompiler(q Query) *Compiler {
	return &Compiler{
		query: q,
		tools: &Tools{
			Params:  params.New(),
			Aliases: aliases.New(),
		},
	}
}

// --- Methods

func (c *Compiler) Expr() *expr.Builder {
	return expr.NewBuilder(c.tools.Params, c.tools.Aliases)
}

func (c *Compiler) Compile(rend expr.Renderer) (*Result, error) {
	w := writer.New(expr.NewParser(rend))

	if err := c.query.Build(c.tools); err != nil {
		return nil, fmt.Errorf("Building time error: %w", err)
	}
	if err := c.query.Render(w); err != nil {
		return nil, fmt.Errorf("Render time error: %w", err)
	}

	return &Result{
		QueryString: w.ToString(),
		Params:      c.tools.Params,
	}, nil
}

type Result struct {
	QueryString string
	Params      core.ParamStore
}
