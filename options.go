package inyorm

import "github.com/laacin/inyorm/internal/core"

// ---- Export internal types

type (
	ColumnType      = core.ColumnType
	AutoPlaceholder = core.AutoPlaceholder
	ColumnWriter    = core.ColumnWriter
)

const (
	TypeColumnDef   = core.ColTypDef
	TypeColumnExpr  = core.ColTypExpr
	TypeColumnAlias = core.ColTypAlias
	TypeColumnBase  = core.ColTypBase
)

type Options struct {
	// AutoPlaceholder
	// Allows selecting in which clauses to automatically insert placeholders.
	// If the provided value is a literal, a placeholder (? | $1) is written instead,
	// and the original value is stored.
	//
	// If nil, the default configuration is:
	//   - Insert = true
	//   - Update = true
	//   - Where  = true
	//   - Having = false
	//   - Join   = false
	AutoPlaceholder *AutoPlaceholder

	// ColumnWriter
	// Allows selecting the default way a column would be written
	//
	// the default configuration is:
	//  - Select  = TypeColumnDef
	//  - Join    = TypeColumnBase
	//  - Where   = TypeColumnExpr
	//  - GroupBy = TypeColumnExpr
	//  - Having  = TypeColumnExpr
	//  - OrderBy = TypeColumnAlias
	ColumnWriter *ColumnWriter

	// ColumnTag
	// Defines the tag that the mapper uses to read and bind values.
	//
	// Default:
	//  - "col"
	ColumnTag string

	// Limit
	// Sets an auto limit clause if the provided value is greater than 0
	Limit int

	// MaxLimit
	// If the provided limit is greater than the maximum, it will be capped at this value
	MaxLimit int
}

// ---- Resolves

func resolveOpts(dialect string, opts **Options) *core.Config {
	ptr := *opts
	if ptr == nil {
		ptr = &Options{}
	}

	core.ResolveAutoPlaceholder(&ptr.AutoPlaceholder)
	core.ResolveColumnWriter(&ptr.ColumnWriter)

	if ptr.ColumnTag == "" {
		ptr.ColumnTag = core.DefaultColumnTag
	}

	return &core.Config{
		Dialect:   dialect,
		AutoPh:    *ptr.AutoPlaceholder,
		ColWrite:  *ptr.ColumnWriter,
		ColumnTag: ptr.ColumnTag,
		MaxLimit:  ptr.MaxLimit,
		Limit:     ptr.Limit,
	}
}
