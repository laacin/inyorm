package column

import "github.com/laacin/inyorm/intr/dialect"

type ColumnBuilder[Col, Case any] struct {
	mainTbl string
	dial    dialect.Dialect
}

// --- Constructors
func (c *ColumnBuilder[Col, Case]) Col(name string, ref ...string) Col {
	tbl := getTbl(c.mainTbl, ref)

	col := dialect.Column{Name: name, Table: tbl}
	return any(&ColumnImpl[Col]{Column: col}).(Col)
}

func (c *ColumnBuilder[Col, Case]) All(ref ...string) Col {
	tbl := getTbl(c.mainTbl, ref)

	col := dialect.Column{Table: tbl, Name: c.dial.Wildcard()}
	return any(&ColumnImpl[Col]{Column: col}).(Col)
}

func (c *ColumnBuilder[Col, Case]) Concat(values ...any) Col {
	col := dialect.Column{Complex: c.dial.ColConcat(values)}
	return any(&ColumnImpl[Col]{Column: col}).(Col)
}

func (c *ColumnBuilder[Col, Case]) Switch(cond any, cas Case) Col {
	cs := any(cas).(dialect.CaseCond)

	col := dialect.Column{Complex: c.dial.ColSwitch(cond, cs)}
	return any(col).(Col)
}

func (c *ColumnBuilder[Col, Case]) Search(cas Case) Col {
	cs := any(cas).(dialect.CaseCond)

	col := dialect.Column{Complex: c.dial.ColSearch(cs)}
	return any(col).(Col)
}

// helpers
func getTbl(main string, ref []string) string {
	if len(ref) > 0 {
		return ref[0]
	}
	return main
}
