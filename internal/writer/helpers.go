package writer

import "github.com/laacin/inyorm/internal/core"

func isAuthPh(ctx core.ClauseType, cfg *core.AutoPlaceholder) bool {
	switch ctx {
	case core.ClsTypInsert:
		return cfg.Insert

	case core.ClsTypUpdate:
		return cfg.Update

	case core.ClsTypWhere:
		return cfg.Where

	case core.ClsTypHaving:
		return cfg.Having

	case core.ClsTypJoin:
		return cfg.Join

	default:
		return false
	}
}

func inferColumn(ctx core.ClauseType, cfg *core.ColumnWriter) core.ColumnType {
	switch ctx {
	case core.ClsTypSelect:
		return cfg.Select

	case core.ClsTypJoin:
		return cfg.Join

	case core.ClsTypWhere:
		return cfg.Where

	case core.ClsTypGroupBy:
		return cfg.GroupBy

	case core.ClsTypHaving:
		return cfg.Having

	case core.ClsTypOrderBy:
		return cfg.OrderBy

	default:
		return core.ColTypExpr
	}
}
