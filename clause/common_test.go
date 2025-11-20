package clause_test

import (
	"reflect"
	"testing"

	"github.com/laacin/inyorm"
	"github.com/laacin/inyorm/clause"
	"github.com/laacin/inyorm/internal/column"
	"github.com/laacin/inyorm/internal/core"
	"github.com/laacin/inyorm/internal/writer"
)

type (
	runner = func(t *testing.T, cls string, vals []any)

	colBuilder = column.ColBuilder[
		inyorm.Column, inyorm.Condition, inyorm.ConditionNext,
		inyorm.CaseSwitch, inyorm.CaseSearch, inyorm.CaseNext,
		inyorm.Identifier, inyorm.Value,
	]

	clsSelect  = clause.Select[inyorm.SelectNext, inyorm.Identifier]
	clsJoin    = clause.Join[inyorm.JoinNext, inyorm.Condition, inyorm.ConditionNext, inyorm.Identifier, inyorm.Value]
	clsWhere   = clause.Where[inyorm.Condition, inyorm.ConditionNext, inyorm.Identifier, inyorm.Value]
	clsGroupBy = clause.GroupBy[inyorm.Identifier]
	clsHaving  = clause.Having[inyorm.Condition, inyorm.ConditionNext, inyorm.Identifier, inyorm.Value]
	clsOrderBy = clause.OrderBy[inyorm.OrderByNext, inyorm.Identifier]
)

func New[T any](cls core.Clause, dialect []string) (T, inyorm.ColumnBuilder, runner) {
	d := ""
	if len(dialect) > 0 {
		d = dialect[0]
	}
	tbl := "users"

	q := writer.NewQuery(tbl, &core.Config{
		Dialect:   d,
		ColumnTag: core.DefaultColumnTag,
		AutoPh:    core.DefaultAutoPlaceholder,
		ColWrite:  core.DefaultColumnWriter,
	})

	clause := cls.(T)
	c := &colBuilder{Table: tbl}

	run := func(t *testing.T, clause string, vals []any) {
		q.SetClauses([]core.Clause{cls}, []core.ClauseType{cls.Name()})
		statement, values := q.Build()

		if statement != clause {
			t.Errorf("\nmismatch result:\nExpect:\n%s\nHave:\n%s\n", clause, statement)
		}

		if !reflect.DeepEqual(vals, values) {
			t.Errorf("\nmissmatch values:\nExpect:\n%#v\nHave:\n%#v\n", vals, values)
		}
	}

	return clause, c, run
}

func NewSelect(dialect ...string) (inyorm.Select, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Select](&clsSelect{}, dialect)
}

func NewJoin(dialect ...string) (inyorm.Join, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Join](&clsJoin{}, dialect)
}

func NewWhere(dialect ...string) (inyorm.Where, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Where](&clsWhere{}, dialect)
}

func NewGroupBy(dialect ...string) (inyorm.GroupBy, inyorm.ColumnBuilder, runner) {
	return New[inyorm.GroupBy](&clsGroupBy{}, dialect)
}

func NewHaving(dialect ...string) (inyorm.Having, inyorm.ColumnBuilder, runner) {
	return New[inyorm.Having](&clsHaving{}, dialect)
}

func NewOrderBy(dialect ...string) (inyorm.OrderBy, inyorm.ColumnBuilder, runner) {
	return New[inyorm.OrderBy](&clsOrderBy{}, dialect)
}
