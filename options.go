package inyorm

import "github.com/laacin/inyorm/internal/core"

// ---- Export internal types

type ColumnType = core.ColumnType

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
	// If nil, the default configuration is:
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

type AutoPlaceholder struct {
	Insert bool
	Update bool
	Where  bool
	Having bool
	Join   bool
}

type ColumnWriter struct {
	Select  ColumnType
	Join    ColumnType
	Where   ColumnType
	GroupBy ColumnType
	Having  ColumnType
	OrderBy ColumnType
}

// ---- Resolves

func resolveOpts(opts **Options) {
	ptr := *opts
	if ptr == nil {
		ptr = &Options{}
	}

	resolveAutoPlaceholder(&ptr.AutoPlaceholder)
	resolveColumnWriter(&ptr.ColumnWriter)

	if ptr.ColumnTag == "" {
		ptr.ColumnTag = "col"
	}

	*opts = ptr
}

func resolveAutoPlaceholder(opt **AutoPlaceholder) {
	ptr := *opt
	if ptr == nil {
		ptr = &AutoPlaceholder{
			Insert: true,
			Update: true,
			Where:  true,
		}
	}
	*opt = ptr
}

func resolveColumnWriter(opt **ColumnWriter) {
	ptr := *opt
	if ptr == nil {
		ptr = &ColumnWriter{}
	}

	dflt := func(provided *ColumnType, dflt ColumnType) {
		if *provided != core.ColTypUnset {
			return
		}
		*provided = dflt
	}

	dflt(&ptr.Select, TypeColumnDef)
	dflt(&ptr.Join, TypeColumnBase)
	dflt(&ptr.Where, TypeColumnExpr)
	dflt(&ptr.GroupBy, TypeColumnExpr)
	dflt(&ptr.Having, TypeColumnExpr)
	dflt(&ptr.OrderBy, TypeColumnAlias)

	*opt = ptr
}
