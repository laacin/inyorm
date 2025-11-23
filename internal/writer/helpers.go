package writer

import "github.com/laacin/inyorm/internal/core"

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
